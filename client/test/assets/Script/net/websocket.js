var websocket = {
    sock: null,
    openHandle: null,
    cloesHandle: null,
    handlers: {},
 
    addHandler: function (messageName,messageDesTool,callback) {
        let tools = {}
        tools.desTool = messageDesTool
        tools.callback = callback
        this.handlers[messageName] = tools
    },
    on_open: function () {
        console.log("client on_open");
        if (this.openHandle)
        {
            this.openHandle()
        }
        /*this.send_data(JSON.stringify({
            stype: "auth",
            ctype: "login",
            data: {
                name: "jianan",
                pwd: 123456
            }
        }));*/
    },

    trigger: function (messageName,data) {
        let handleTool = this.handlers[messageName]
        let message = data.getMessagecontent().unpack(handleTool.desTool,messageName)
        handleTool.callback(message)
    },
    
    on_message: function (event) {
        console.log("client rcv data");
        let data = proto.ClientReply.deserializeBinary(event.data);
        console.log("data",data);
        let messageName = data.getMessagename()
        this.trigger(messageName,data)
        //console.log("client rcv data=" + event.data);
    },
 
    on_close: function () {
        console.log("client on_close")
        this.close();
        if (this.cloesHandle)
        {
            this.cloesHandle()
        }
    },
 
    on_error: function () {
        console.log("client on_error")
        this.close();
    },
    
    close: function () {
        if(this.sock){
            this.sock.close();
            this.sock = null;
        }
    },
 
    connect: function (url) {
        this.sock = new WebSocket(url);
        this.sock.binaryType = "arraybuffer";
        this.sock.onopen = this.on_open.bind(this);
        this.sock.onmessage = this.on_message.bind(this);
        this.sock.onclose = this.on_close.bind(this);
        this.sock.onerror = this.on_error.bind(this);
    },
 
    send_data: function (data) {
        this.sock.send(data);
    },
 
    initHandle: function (openHandle,cloesHandle) {
        this.openHandle = openHandle;
        this.cloesHandle = cloesHandle;
    }
}
 
module.exports = websocket;