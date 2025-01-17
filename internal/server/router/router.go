package router

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/AndrXxX/goph-keeper/internal/enums/datatypes"
	"github.com/AndrXxX/goph-keeper/internal/server/api/middlewares"
	"github.com/AndrXxX/goph-keeper/internal/server/config"
	"github.com/AndrXxX/goph-keeper/internal/server/controllers"
	"github.com/AndrXxX/goph-keeper/internal/server/entities"
	"github.com/AndrXxX/goph-keeper/internal/server/services/entityconvertors"
	"github.com/AndrXxX/goph-keeper/internal/server/services/valueconvertors"
	"github.com/AndrXxX/goph-keeper/pkg/hashgenerator"
	"github.com/AndrXxX/goph-keeper/pkg/requestjsonentity"
	"github.com/AndrXxX/goph-keeper/pkg/token"
)

type router struct {
	config  appConfig
	storage Storage
}

func New(c *config.Config, storage Storage) *router {
	return &router{config: appConfig{c}, storage: storage}
}

func (mr *router) RegisterApi() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	mr.registerAPI(r)
	return r
}

func (mr *router) registerAPI(r *chi.Mux) {
	hg := hashgenerator.Factory().SHA256(mr.config.c.PasswordKey)
	ts := token.New(mr.config.c.AuthKey, time.Duration(mr.config.c.AuthKeyExpired)*time.Second)

	r.Use(middlewares.RequestLogger().Handler)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.CompressGzip().Handler)
		ac := controllers.AuthController{US: mr.storage.US, HG: hg, TS: ts, UF: &requestjsonentity.Fetcher[entities.User]{}, KeyPath: mr.config.c.PublicKeyPath}
		r.Post("/api/user/register", ac.Register)
		r.Post("/api/user/login", ac.Login)
	})

	r.Group(func(r chi.Router) {
		ecf := entityconvertors.Factory{}
		vcf := valueconvertors.Factory{}

		r.Use(middlewares.IsAuthorized(ts).Handler)
		r.Use(middlewares.CompressGzip().Handler)
		lpc := controllers.ItemsController[entities.PasswordItem]{
			Type:      datatypes.Passwords,
			Fetcher:   &requestjsonentity.Fetcher[entities.PasswordItem]{},
			Storage:   mr.storage.IS,
			Convertor: ecf.Password(vcf.Password()),
		}
		r.Post(fmt.Sprintf("/api/updates/%s", datatypes.Passwords), lpc.StoreUpdates)
		r.Get(fmt.Sprintf("/api/updates/%s", datatypes.Passwords), lpc.FetchUpdates)

		tc := controllers.ItemsController[entities.NoteItem]{
			Type:      datatypes.Notes,
			Fetcher:   &requestjsonentity.Fetcher[entities.NoteItem]{},
			Storage:   mr.storage.IS,
			Convertor: ecf.Note(vcf.Note()),
		}
		r.Post(fmt.Sprintf("/api/updates/%s", datatypes.Notes), tc.StoreUpdates)
		r.Get(fmt.Sprintf("/api/updates/%s", datatypes.Notes), tc.FetchUpdates)

		bcc := controllers.ItemsController[entities.BankCardItem]{
			Type:      datatypes.BankCards,
			Fetcher:   &requestjsonentity.Fetcher[entities.BankCardItem]{},
			Storage:   mr.storage.IS,
			Convertor: ecf.BankCard(vcf.BankCard()),
		}
		r.Post(fmt.Sprintf("/api/updates/%s", datatypes.BankCards), bcc.StoreUpdates)
		r.Get(fmt.Sprintf("/api/updates/%s", datatypes.BankCards), bcc.FetchUpdates)

		bc := controllers.ItemsController[entities.FileItem]{
			Type:      datatypes.Files,
			Fetcher:   &requestjsonentity.Fetcher[entities.FileItem]{},
			Storage:   mr.storage.IS,
			Convertor: ecf.File(vcf.File()),
		}
		r.Get(fmt.Sprintf("/api/updates/%s", datatypes.Files), bc.FetchUpdates)

		fc := controllers.FilesController{
			Storage:   mr.storage.IS,
			FS:        mr.storage.FS,
			FF:        &requestjsonentity.Fetcher[entities.FileItem]{},
			Convertor: ecf.File(vcf.File()),
		}
		r.Post("/api/files/update", fc.Update)
		r.Post("/api/files/upload/{id}/", fc.Upload)
		r.Get("/api/files/download/{id}/", fc.Download)
	})
}
