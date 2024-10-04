package errgroup

import (
	"context"
	"sync"
)

/*
Одно из заданий было написать свою реализацию errgroup.
Чтобы было честно я глянул на реализацию один раз, и теперь попробую без подглядок написать свою
Попытаемся сформулировать требования.
1. Должен быть метод Go, который запускает функцию в горутине, если только мы не получили ошибку ранее
2. Как только одна из функций вернула ошибку, она должна быть сразу же возвращена без ожидания завершения остальных
3. Должен быть метод Wait, который блокирует выполнение до тех пор, пока все либо все функции не завершаться, либо не
возникнет первая ошибка, повторный вызов метода возвращает первую полученную ошибку или nil
4. ErrGroup должна поддерживать только режим с контекстом, чтобы отменять все горутины при ошибке
5. Поддерживается режим только с отменой контекста, так как это кастомная ошибка - я могу себе это позволить
*/

type errGroup struct {
	wg       sync.WaitGroup
	cancel   context.CancelFunc
	errCh    chan error
	firstErr error
	done     chan struct{}
}

type ErrGroup interface {
	Go(errFn func() error)
	Wait() error
}

func NewErrGroup(ctx context.Context) (ErrGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &errGroup{
		errCh:  make(chan error),
		done:   make(chan struct{}),
		cancel: cancel,
	}, ctx
}

func (e *errGroup) Go(errFn func() error) {
	select {
	case <-e.done:
		return
	default:
	}
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		err := errFn()
		if err != nil {
			select {
			case <-e.done:
				return
			case e.errCh <- err:
				e.cancel()
				close(e.done)
			}
		}
	}()
}

func (e *errGroup) Wait() error {
	go func() {
		e.wg.Wait()
		select {
		case <-e.done:
			return
		case e.errCh <- nil:
			e.cancel()
			close(e.done)
		}
	}()
	select {
	case <-e.done:
		return e.firstErr
	case e.firstErr = <-e.errCh:
		return e.firstErr
	}
}

/*
Вот и написана моя реализация errgroup, понимаю, что она мягко говоря не идеальна, но так же интереснее, поэтому жду
твоего фидбэка, чтобы понять, что можно улучшить. Тесты с ней прогонял, с флагом -race в том числе - проблем не обнаружил
*/
