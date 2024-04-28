package main

import (
	"context"
	"math/rand"
	"time"

	plugin "github.com/sandrolain/event-runner/src/plugin/bootstrap"
	"github.com/sandrolain/event-runner/src/plugin/proto"
)

func main() {
	plugin.Start(plugin.StartOptions{
		Service: &server{},
		Callback: func() error {

			time.Sleep(2 * time.Second)

			plugin.SetReady()

			return nil
		},
	})
}

type server struct {
	proto.UnimplementedPluginServiceServer
}

func (s *server) Status(ctx context.Context, in *proto.StatusReq) (*proto.StatusRes, error) {
	return plugin.GetStatusResponse(), nil
}

func (s *server) Shutdown(ctx context.Context, in *proto.ShutdownReq) (*proto.ShutdownRes, error) {
	return plugin.Shutdown(in.Wait), nil
}

func (s *server) Command(ctx context.Context, in *proto.CommandReq) (*proto.CommandRes, error) {
	if in.Command == "pizza" {
		s := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s)
		i := r.Intn(len(pizzas))
		d := pizzas[i]
		return plugin.SuccessResult(in, d)
	}
	return plugin.NotFoundResult(in)
}

var pizzas = []string{
	"Margherita", "Marinara", "Quattro Stagioni", "Carbonara", "Frutti di Mare",
	"Quattro Formaggi", "Crudo", "Napoletana", "Pugliese", "Montanara",
	"Emiliana", "Romana", "Fattoria", "Capricciosa", "Siciliana",
	"Ortolana", "Prosciutto", "Funghi", "Calzone", "Pepperoni",
}
