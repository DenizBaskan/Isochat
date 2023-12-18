import Navbar from './Navbar'
import useWebSocket, { ReadyState } from 'react-use-websocket'
import { useState , useEffect, useRef, useLayoutEffect } from 'react'
import { domain, ws_url } from '../Globals'
import axios from 'axios'

const Status = {
    SendMessage: 0
}

export default function Chat() {
    document.title = "Chat"
    
    const [friendRequestUsername, setFriendRequestUsername] = useState("")
    const [message, setMessage] = useState("")
    const [recipientID, setRecpientID] = useState("")
    const [error, setError] = useState("")

    const [messages, setMessages] = useState([])

    const [friends, setFriends] = useState(new Map())
    const [pending, setPending] = useState(new Map())
    const [incoming, setIncoming] = useState(new Map())

    const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(
        ws_url, {
            share: false,
            shouldReconnect: () => true,
        }
    )

    useEffect(() => {
        var msg = lastJsonMessage
        
        if (msg != null) {
            if (!msg.success) {
                setError(msg.reason)
            } else {
                setError("")

                if (msg.status == Status.SendMessage) {
                    var newMessages = [...messages]
                    newMessages.push(msg.data.updates.message)
                    setMessages(newMessages)
                }
            }
        }
    }, [lastJsonMessage])

    useEffect(() => {
        axios.get(domain + "/friends", {withCredentials : true}).then((res) => {
            if (res.data) {
                if (res.data.friends) {
                    var newFriends = new Map()

                    res.data.friends.forEach((f) => {
                        newFriends.set(f.id, {
                            user_id: f.user_id,
                            username: f.username
                        })
                    })

                    setFriends(newFriends)
                }
                
                if (res.data.pending) {
                    var newPending = new Map()

                    res.data.pending.forEach((f) => {
                        newPending.set(f.id, {
                            user_id: f.user_id,
                            username: f.username
                        })
                    })

                    setPending(newPending)
                }
                
                if (res.data.incoming) {
                    var newIncoming = new Map()

                    res.data.incoming.forEach((f) => {
                        newIncoming.set(f.id, {
                            user_id: f.user_id,
                            username: f.username
                        })
                    })

                    setIncoming(newIncoming)
                }
            }
        }).catch((e) => {
            if (e.response) {
                if (e.response.status == 400) {
                    return setError(e.response.data["reason"])
                }
            }

            return setError("An unexpected error occurred")
        })
    }, [])

    useEffect(() => {
        if (recipientID != "")
            axios.get(domain + "/messages/" + recipientID, {withCredentials : true}).then((res) => {
                if (res.data) {
                    setMessages(res.data)
                }
            }).catch((e) => {
                if (e.response) {
                    if (e.response.status == 400) {
                        return setError(e.response.data["reason"])
                    }
                }
    
                return setError("An unexpected error occurred")
            })
    }, [recipientID])

    const sendMessage = (event) => {
        event.preventDefault()
        
        if (readyState === ReadyState.OPEN) {
            sendJsonMessage({
                status: Status.SendMessage,
                data: {
                    recipient_id: recipientID,
                    message: message
                }
            })

            setMessage("")
        } else {
            setError("An unexpected error occured")
        }
    }

    const sendRequest = (event) => {
        event.preventDefault()

        axios.post(domain + "/friend/request", {
            username: friendRequestUsername
        }, {withCredentials : true}).then(() => {
            window.location.reload()
        }).catch((e) => {
            if (e.response) {
                if (e.response.status == 400) {
                    return setError(e.response.data["reason"])
                }
            }

            return setError("An unexpected error occurred")
        })
    }

    const acceptRequest = (id) => {
        axios.post(domain + "/friend/request/accept", {
            friend_id: id
        }, {withCredentials : true}).then(() => {
            window.location.reload()
        }).catch((e) => {
            if (e.response) {
                if (e.response.status == 400) {
                    return setError(e.response.data["reason"])
                }
            }

            return setError("An unexpected error occurred")
        })
    }

    const declineRequest = (id) => {
        axios.post(domain + "/friend/request/decline", {
            friend_id: id
        }, {withCredentials : true}).then(() => {
            window.location.reload()
        }).catch((e) => {
            if (e.response) {
                if (e.response.status == 400) {
                    return setError(e.response.data["reason"])
                }
            }

            return setError("An unexpected error occurred")
        })
    }

    const removeFriend = (id) => {
        axios.delete(domain + "/friend/" + id, {withCredentials : true}).then(() => {
            window.location.reload()
        }).catch((e) => {
            if (e.response) {
                if (e.response.status == 400) {
                    return setError(e.response.data["reason"])
                }
            }

            return setError("An unexpected error occurred")
        })
    }

    const removeFriendRequest = (id) => {
        axios.delete(domain + "/friend/request/" + id, {withCredentials : true}).then(() => {
            window.location.reload()
        }).catch((e) => {
            if (e.response) {
                if (e.response.status == 400) {
                    return setError(e.response.data["reason"])
                }
            }

            return setError("An unexpected error occurred")
        })
    }

    return (
        <>
            <Navbar />

            <p className="text-center m-3 text-danger">{error}</p>
            
            <form className="text-center" onSubmit={sendRequest}>
                <p>Send a friend request</p>
                <input className="mx-auto" placeholder="Username" value={friendRequestUsername} onChange={e => setFriendRequestUsername(e.target.value)}/>
                <input className=" btn btn-primary m-2" type="submit" value="Send"/>
            </form>

            <div className="row m-2">
                <div className="m-5" id="friend-box">
                    <div className="p-3">
                        <h5>Friends</h5>
                        {Array.from(friends.entries()).map((f) => {
                            const [k, v] = f
                            return (<p><a href="javascript:void(0);" onClick={() => setRecpientID(v.user_id)}>@{v.username}</a> <a href="javascript:void(0);" onClick={() => removeFriend(k)}>Remove</a></p>)
                        })}
                        <h5>Incoming</h5>
                        {Array.from(incoming.entries()).map((f) => {
                            const [k, v] = f
                            return (<p>@{v.username} <a href="javascript:void(0);" onClick={() => acceptRequest(k)}>Accept</a> <a href="javascript:void(0);" onClick={() => declineRequest(k)}>Decline</a></p>)
                        })}
                        <h5>Pending</h5>
                        {Array.from(pending.entries()).map((f) => {
                            const [k, v] = f
                            return (<p>@{v.username} <a href="javascript:void(0);" onClick={() => removeFriendRequest(k)}>Remove</a></p>)
                        })}
                    </div>
                </div>

                <div className="justify-content-center m-5" id="chat-box">
                    <div className="p-4">
                        {messages != null && messages.map(function(m, i) {
                            return <p className={m.is_sender ? "message-recieved": ""}>{m.is_sender ? "" : "@" + m.sender_username + ": "}{m.data}</p>
                        })}
                    </div>
                </div>
            </div>

            <form className="text-center" onSubmit={sendMessage}>
                <input className="m-4 mx-auto m-3" placeholder="Message" value={message} onChange={e => setMessage(e.target.value)}/>
                <input className="m-2 btn btn-primary" type="submit" value="Send"/>
            </form>
        </>
    )
}
