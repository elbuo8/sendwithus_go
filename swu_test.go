package swu

import (
	"log"
	"os"
	"testing"
)

func TestNewSWU(t *testing.T) {
	api := New("key")
	if api == nil {
		t.Error("New should not return nil")
	}
}

func TestTemplates(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	a, err := api.Emails()
	if err != nil {
		t.Error(err)
	}
	log.Print(a)
}

func TestGetTemplate(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	a, err := api.GetTemplate("tem_bRKXvNLAXTG8EGxhut3gCe")
	if err != nil {
		t.Error(err)
	}
	log.Print(a)
}

func TestGetTemplateVersion(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	a, err := api.GetTemplateVersion("tem_bRKXvNLAXTG8EGxhut3gCe", "ver_Hh35dZhnffghidEy6VeHKL")
	if err != nil {
		t.Error(err)
	}
	log.Print(a)
}

func TestUpdateTemplateVersion(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	a, err := api.UpdateTemplateVersion("tem_bRKXvNLAXTG8EGxhut3gCe", "ver_Hh35dZhnffghidEy6VeHKL",
		&SWUVersion{
			Name:    "Test",
			Subject: "Test",
			Text:    "test",
		})
	if err != nil {
		t.Error(err)
	}
	log.Print(a)
}

func TestCreateTemplate(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	a, err := api.CreateTemplate(&SWUVersion{
		Name:    "test",
		Subject: "test",
		Text:    "ALOHA",
	})
	if err != nil {
		t.Error(err)
	}
	log.Print(a)
}

func TestCreateTemplateVersion(t *testing.T) {
	api := New(os.Getenv("SWU_KEY"))
	a, err := api.CreateTemplateVersion("tem_nXAPFGXQXFKcibJHdm9PZ9", &SWUVersion{
		Name:    "test",
		Subject: "test",
		Text:    "ALOHA1",
	})
	if err != nil {
		t.Error(err)
	}
	log.Print(a)
}
