package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}
type NotificationService struct {
	// not need to implement
	NotEmptyStruct bool
}
type PenaltyService struct {
	// not need to implement
	Сoefficient int
}

type Container struct {
	constructors map[string]reflect.Value
	singletons   map[string]*reflect.Value
}

func NewContainer() *Container {
	constructors := make(map[string]reflect.Value)
	singletons := make(map[string]*reflect.Value)
	return &Container{
		constructors: constructors,
		singletons:   singletons,
	}
}

func (c *Container) RegisterSingletonType(name string, constructor any) {
	if _, exists := c.constructors[name]; !exists {
		c.constructors[name] = reflect.ValueOf(constructor)
		c.singletons[name] = nil
	}
}

func (c *Container) RegisterType(name string, constructor any) {
	if _, exists := c.constructors[name]; !exists {
		c.constructors[name] = reflect.ValueOf(constructor)
	}
}

func (c *Container) Resolve(name string, args ...any) (any, error) {
	singleton, existsSingleton := c.singletons[name]
	if existsSingleton && singleton != nil {
		return singleton.Interface(), nil
	}

	if constructor, exists := c.constructors[name]; exists {
		cType := constructor.Type()
		n := cType.NumIn()
		a := make([]reflect.Value, n)
		for i := range n {
			if reflect.TypeOf(args[i]) == cType.In(i) {
				a[i] = reflect.ValueOf(args[i])
			}
		}

		instance := constructor.Call(a)[0]
		if existsSingleton && singleton == nil {
			c.singletons[name] = &instance
		}

		return instance.Interface(), nil
	}

	return nil, fmt.Errorf("constructor %s not found", name)
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})
	container.RegisterSingletonType("NotificationService", func() interface{} {
		return &NotificationService{}
	})
	container.RegisterType("PenaltyService", func(arg int) interface{} {
		return &PenaltyService{
			Сoefficient: arg,
		}
	})

	penaltyService1, err := container.Resolve("PenaltyService", 12)
	assert.NoError(t, err)
	penaltyService2, err := container.Resolve("PenaltyService", 22)
	assert.NoError(t, err)

	p1 := penaltyService1.(*PenaltyService)
	assert.Equal(t, 12, p1.Сoefficient)
	p2 := penaltyService2.(*PenaltyService)
	assert.Equal(t, 22, p2.Сoefficient)
	assert.False(t, p1 == p2)

	notificationService1, err := container.Resolve("NotificationService")
	assert.NoError(t, err)
	notificationService2, err := container.Resolve("NotificationService")
	assert.NoError(t, err)

	n1 := notificationService1.(*NotificationService)
	n2 := notificationService2.(*NotificationService)
	assert.True(t, n1 == n2)

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
