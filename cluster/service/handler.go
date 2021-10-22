package service

import (
	"github.com/aceaura/libra/cluster"
)

type Handler struct {
	codec cluster.Codec
	name  string
	code  uint64
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(r *Router) {
	
}

func (h *Handler) init() {

}

func (h *Handler) WithCodec(codec cluster.Codec) *Handler {
	return h
}

func (h *Handler) WithName(name string) *Handler {
	return h
}

func (h *Handler) WithCode(code uint64) *Handler {
	return h
}
