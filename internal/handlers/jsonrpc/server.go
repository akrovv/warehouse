package jsonrpc

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/akrovv/warehouse/pkg/logger"
)

type server struct {
	server *rpc.Server
}

type HTTPConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HTTPConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HTTPConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HTTPConn) Close() error                      { return nil }

func NewServer(productService ProductService, warehouseService WarehouseService, logger logger.Logger) (*server, error) {
	r := rpc.NewServer()

	if err := r.RegisterName("Products", NewProductHandler(productService, logger)); err != nil {
		return nil, err
	}

	if err := r.RegisterName("Warehouses", NewWarehouseHandler(warehouseService, logger)); err != nil {
		return nil, err
	}

	return &server{
		server: r,
	}, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serverCodec := jsonrpc.NewServerCodec(&HTTPConn{
		in:  r.Body,
		out: w,
	})

	w.Header().Set("Content-type", "application/json")
	err := s.server.ServeRequest(serverCodec)
	if err != nil {
		http.Error(w, `{"error":"cant serve request"}`, http.StatusInternalServerError)
	}
}

func (s *server) Run(port string) error {
	http.Handle("/", s)
	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
