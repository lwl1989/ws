package component

//import (
//    "context"
//    "fmt"
//    "go.etcd.io/etcd/clientv3"
//    "time"
//)
//
////创建租约注册服务
//type ServiceReg struct {
//    client          *clientv3.Client
//    lease           clientv3.Lease
//    leaseResp       *clientv3.LeaseGrantResponse
//    cancel          func()
//    keepAliveChan   <-chan *clientv3.LeaseKeepAliveResponse
//    key             string
//}
//
//func NewServiceReg(addr []string, timeNum, timeout int64,) (*ServiceReg, error) {
//    conf := clientv3.Config{
//        Endpoints:   addr,
//        DialTimeout: time.Duration(timeout * int64(time.Second)),
//    }
//
//    client, err := clientv3.New(conf)
//    if err != nil {
//        return nil, err
//    }
//
//    ser := &ServiceReg{
//        client: client,
//    }
//
//    if err := ser.setLease(timeNum); err != nil {
//        return nil, err
//    }
//    go ser.ListenLeaseRespChan()
//    return ser, nil
//}
//
////设置租约
//func (reg *ServiceReg) setLease(timeNum int64) error {
//    lease := clientv3.NewLease(reg.client)
//
//    //设置租约时间
//    leaseResp, err := lease.Grant(context.TODO(), timeNum)
//    if err != nil {
//        return err
//    }
//
//    //设置续租
//    ctx, cancelFunc := context.WithCancel(context.TODO())
//    leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)
//
//    if err != nil {
//        return err
//    }
//
//    reg.lease = lease
//    reg.leaseResp = leaseResp
//    reg.cancel = cancelFunc
//    reg.keepAliveChan = leaseRespChan
//    return nil
//}
//
////监听 续租情况
//func (reg *ServiceReg) ListenLeaseRespChan() {
//    for {
//        select {
//        case leaseKeepResp := <-reg.keepAliveChan:
//            if leaseKeepResp == nil {
//                fmt.Println("已经关闭续租功能")
//                return
//            } else {
//                fmt.Println("续租成功")
//            }
//        }
//    }
//}
//
////通过租约 注册服务
//func (reg *ServiceReg) PutService(key, val string) error {
//    kv := clientv3.NewKV(reg.client)
//    _, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(reg.leaseResp.ID))
//    return err
//}
//
//
////撤销租约
//func (reg *ServiceReg) RevokeLease() error {
//    reg.cancel()
//    time.Sleep(2 * time.Second)
//    _, err := reg.lease.Revoke(context.TODO(), reg.leaseResp.ID)
//    return err
//}
