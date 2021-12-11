import React, {useContext } from 'react';
import { ChannelList } from './ChannelList';
import './chat.scss';
import { MessagesPanel } from './MessagesPanel';
import { Modal } from 'react-bootstrap';
import 'bootstrap/dist/css/bootstrap.min.css';
import { Navigate } from 'react-router-dom';
import SelectSearch from 'react-select-search';
import AsyncSelect from 'react-select/async'
import socketClient from "socket.io-client";
const SERVER = "http://127.0.0.1:8080";

export class Chat extends React.Component {
    state = {
        channels: [],
        socket: null,
        channel: null,
        showModal: false,
        options: [],
        listed_user: [],
        new_name_room: null,
        username: null,
        id: null,
        auth: false
    }
    socket;
    conn;
    componentDidMount() {
        const loggedInUser = localStorage.getItem("username");
        const loggedInID = localStorage.getItem("id");
        if (loggedInUser && loggedInID) {
            this.state.auth = true;
            this.state.username = loggedInUser
            this.state.id = loggedInID

            this.loadChannels();
            this.configureSocket();
        }else{
            this.state.auth = false;
        }
    }

    configureSocket = () => {
        if (window["WebSocket"]) {
            const url_string = window.location.href;
            const url = new URL(url_string);
            const roomId = url.searchParams.get("roomid");
            this.conn = new WebSocket("ws://localhost:8080/ws?roomid=1&username="+this.state.username);
            this.conn.onopen = () => {
                console.log('WebSocket Client Connected');
            };
            this.conn.onclose = function (evt) {
                console.log("connection closed")
            };
            this.conn.onmessage = (evt) => {
                let data = JSON.parse(evt.data);
                let channels = this.state.channels;
                let list_users = this.state.options;
                
                if(data && 'flag' in data){
                    if(data.flag == "message"){
                        channels.forEach(c => {
                            if (c.id === data.channel_id) {
                                if (!c.messages) {
                                    c.messages = [data.messagec];
                                } else {
                                    c.messages.push(data.messagec);
                                }
                            }
                        });
                        this.setState({ channels });
                    }else if(data.flag == "channel-join"){
                        channels.forEach(c => {
                            if (c.id === data.channel_id) {
                                c.participants = data.participants;
                                c.messages = data.messageb
                            }
                        });
                        this.setState({ channels });
                    }else if(data.flag == "get-channel"){
                        this.setState({ channels: data.subcriptions });
                        this.setState({ showModal: false });
                    }else {
                        if(data && 'users' in data){
                            let list_opt = [];
                            for(let i in data.users){
                                let opt = {
                                    label: data.users[i].username,
                                    value: data.users[i].id
                                }
                                list_opt.push(opt)
                            }
                            this.state.options = list_opt;
                        }
                    }
                }
            };
        } else {
            console.log("Your browser does not support WebSockets.");
        }
    }

    loadChannels = async () => {
        fetch('http://localhost:8080/getChannels?username='+this.state.username).then(async response => {
            let data = await response.json();
            this.setState({ channels: data });
        })
    }

    handleChannelSelect = id => {
        let channel = this.state.channels.find(c => {
            return c.id === id;
        });
        this.setState({ channel });
        this.conn.send(JSON.stringify({ flag: "channel-join", channel_id: channel.id, room: {name: channel.room}}));
    }

    handleSendMessage = (channel_id, text) => {
        let channel = this.state.channels.find(c => {
            return c.id === channel_id;
        });
        //ID is userid
        this.conn.send(JSON.stringify({ flag: "message", id: parseInt(this.state.id), channel_id: channel_id, room: {name: channel.room}, messagec: {content: text, username: this.state.username}}));
    }

    showModal = () => {
        this.setState({ showModal: true });
    };

    hideModal = () => {
        this.setState({ showModal: false });
    };

    loadOptions = (inputValue, callback) => {
        this.conn.send(JSON.stringify({ flag: "find-user", id: parseInt(this.state.id), text: inputValue}));
        new Promise((resolve) => {
            setTimeout(() => {
              resolve(callback(this.state.options));
            }, 1000);
        });
    }

    createRoom = () => {
        for(let i in this.state.listed_user){
            this.state.listed_user[i].id = this.state.listed_user[i].value
        }
        this.conn.send(JSON.stringify({ flag: "create-room", id: parseInt(this.state.id), users: this.state.listed_user, room: {Public: {rn: this.state.new_name_room}}}));
    }

    handleInput = e => {
        this.setState({ new_name_room: e.target.value });
    }

    render() {
        const loggedInUser = localStorage.getItem("username");
        const loggedInID = localStorage.getItem("id");

        if (loggedInUser && loggedInID) {
            return (
                <>
                <Modal
                    size="lg"
                    show={this.state.showModal}
                    onHide={() => this.state.showModal = false}
                    aria-labelledby="example-modal-sizes-title-lg"
                >
                    <Modal.Header closeButton>
                    <Modal.Title id="example-modal-sizes-title-lg">
                        Create Room
                    </Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                    <input className="input-create" type="text" onChange={this.handleInput} value={this.state.new_name_room} placeholder="Enter Room Name" />
                    <AsyncSelect
                        loadOptions={this.loadOptions}
                        isMulti
                        onChange={opt => this.state.listed_user = opt}
                    />
                    <div className="create">
                    <button className="button-cancel" onClick={this.hideModal}>Cancel</button>
                    <button className="button-create" onClick={this.createRoom}>Create</button>
                    </div>
                    </Modal.Body>
                </Modal>
                <div className='chat-app'>
                    <ChannelList channels={this.state.channels} onSelectChannel={this.handleChannelSelect} onShowModal={this.showModal}/>
                    <MessagesPanel onSendMessage={this.handleSendMessage} channel={this.state.channel} />
                </div>
                </>
            );
        }else{
            return (
                <Navigate to="/" />
            )
        }
    }
}