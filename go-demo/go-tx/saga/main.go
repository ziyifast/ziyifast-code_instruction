package main

import (
	"fmt"
	"log"
	"time"
)

// 定义一个简单的服务接口
type Service interface {
	DoSomething() error
	Compensate() error
}

// 模拟服务A
type ServiceA struct {
	data string
}

func (s *ServiceA) DoSomething() error {
	fmt.Println("Service A: Doing something...")
	// 模拟成功的情况
	s.data = "Service A data"
	return nil
}

func (s *ServiceA) Compensate() error {
	fmt.Println("Service A: Compensating...")
	s.data = "" // 撤销操作
	return nil
}

// 模拟服务B
type ServiceB struct {
	data int
}

func (s *ServiceB) DoSomething() error {
	fmt.Println("Service B: Doing something...")
	// 模拟成功的情况
	s.data = 123
	return nil
}

func (s *ServiceB) Compensate() error {
	fmt.Println("Service B: Compensating...")
	s.data = 0 // 撤销操作
	return nil
}

// Saga orchestrator
type SagaOrchestrator struct {
	services  []Service
	completed []bool
}

func NewSagaOrchestrator(services []Service) *SagaOrchestrator {
	return &SagaOrchestrator{
		services:  services,
		completed: make([]bool, len(services)),
	}
}

func (s *SagaOrchestrator) Run() error {
	for i, service := range s.services {
		err := service.DoSomething()
		if err != nil {
			log.Printf("Service %d failed: %v\n", i, err)
			return s.compensate(i)
		}
		s.completed[i] = true
	}
	return nil
}

func (s *SagaOrchestrator) compensate(failedIndex int) error {
	//任务失败，执行对应服务的补偿措施
	fmt.Println("Starting compensation...")
	for i := failedIndex; i >= 0; i-- {
		if s.completed[i] {
			err := s.services[i].Compensate()
			if err != nil {
				log.Printf("Compensation for service %d failed: %v\n", i, err)
				// 这里可以考虑重试补偿操作，或者记录日志并人工介入
				return err
			}
			s.completed[i] = false
		}
	}
	return nil
}

func main() {
	serviceA := &ServiceA{}
	serviceB := &ServiceB{}

	saga := NewSagaOrchestrator([]Service{serviceA, serviceB})

	err := saga.Run()
	if err != nil {
		log.Fatalf("Saga failed: %v\n", err)
	} else {
		fmt.Println("Saga completed successfully!")
	}

	time.Sleep(time.Second) // 模拟等待
}
