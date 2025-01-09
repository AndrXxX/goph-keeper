package storageadapters

import (
	"github.com/AndrXxX/goph-keeper/internal/client/interfaces"
)

type ORMAdapter[E any, M any] struct {
	Storage         interfaces.Storage[M]
	ORMConvertor    convertor[M, E]
	EntityConvertor convertor[E, M]
}

func (a *ORMAdapter[E, M]) Find(e *E) *E {
	return a.ORMConvertor.Convert(a.Storage.Find(a.EntityConvertor.Convert(e)))
}

func (a *ORMAdapter[E, M]) Create(e *E) (*E, error) {
	m, err := a.Storage.Create(a.EntityConvertor.Convert(e))
	if err != nil {
		return nil, err
	}
	return a.ORMConvertor.Convert(m), nil
}

func (a *ORMAdapter[E, M]) Update(e *E) error {
	return a.Storage.Update(a.EntityConvertor.Convert(e))
}

func (a *ORMAdapter[E, M]) FindAll(e *E) []E {
	mList := a.Storage.FindAll(a.EntityConvertor.Convert(e))
	list := make([]E, len(mList))
	for i := range mList {
		list[i] = *a.ORMConvertor.Convert(&mList[i])
	}
	return list
}
