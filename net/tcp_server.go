package net

type TcpServer interface {
	Start() error                                                                                                 //启动网络服务器
	Stop()                                                                                                        //停止网络服务
	AddAcceptor(hand ListenerHandler, tag Tag, addr string, pause bool, ops ...TcpAceptorOptions) (string, error) //添加监听器
	RemoveAcceptor(tag Tag)                                                                                       //移除指定监听器,并清理其连接
	DisConnectByTag(tag Tag, pause bool)                                                                          //断开指定监听器的连接,并设置暂停接受连接
	DisConnectAll(pause bool)                                                                                     //断开所有连接,并设置暂停接受连接
	PauseAcceptByTag(tag Tag)                                                                                     //暂停指定监听器接受连接
	PauseAccept()                                                                                                 //暂停所有监听器接受连接
	ResumeAcceptByTag(tag Tag)                                                                                    //恢复指定监听器接受连接
	ResumeAccept()                                                                                                //恢复所有监听器接受连接
}
