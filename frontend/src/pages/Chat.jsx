import Navbar from './Navbar'
import useWebSocket, { ReadyState } from "react-use-websocket"
import { useState , useEffect } from 'react'
import { domain, ws_url } from '../Globals'

const Status = {
    SendMessage: 0,
    GetMessages: 1,
    GetFriends: 2,
    SendFriendRequest: 3,
    AcceptFriendRequest: 4,
    DeclineFriendRequest: 5,
    RemoveFriend: 6,
    RemoveFriendRequest: 7
}

export default function Chat() {
    document.title = "Chat"
    
    const [friendRequestUsername, setFriendRequestUsername] = useState("")
    const [message, setMessage] = useState("")
    const [recipientID, setRecpientID] = useState("")
    const [error, setError] = useState("")
    const [friends, setFriends] = useState({})
    
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
                
                if (msg.status == Status.GetFriends) {
                    setFriends(msg.data)
                }
            }
        }
    }, [lastJsonMessage])

    useEffect(() => {
        if (readyState === ReadyState.OPEN) {
            sendJsonMessage({
                status: Status.GetFriends
            })
        }
    }, [readyState])

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
        }

        setMessage("")
    }

    const sendRequest = (event) => {
        event.preventDefault()
        
        if (readyState === ReadyState.OPEN) {
            sendJsonMessage({
                status: Status.SendFriendRequest,
                data: {
                    username: friendRequestUsername
                }
            })
        }

        setMessage("")
    }

    const acceptRequest = (id) => {
        if (readyState === ReadyState.OPEN) {
            sendJsonMessage({
                status: Status.AcceptFriendRequest,
                data: {
                    friend_id: id
                }
            })
        }
    }

    const declineRequest = (id) => {
        if (readyState === ReadyState.OPEN) {
            sendJsonMessage({
                status: Status.DeclineFriendRequest,
                data: {
                    friend_id: id
                }
            })
        }
    }

    const removeFriend = (id) => {
        if (readyState === ReadyState.OPEN) {
            sendJsonMessage({
                status: Status.RemoveFriend,
                data: {
                    friend_id: id
                }
            })
        }
    }

    const removeFriendRequest = (id) => {
        if (readyState === ReadyState.OPEN) {
            sendJsonMessage({
                status: Status.RemoveFriendRequest,
                data: {
                    friend_id: id
                }
            })
        }
    }

    return (
        <>
            <Navbar />

            <p className="text-center m-3 text-danger">{error}</p>
            
            <form className="text-center" onSubmit={sendRequest}>
                <p>Send a friend request</p>
                <input className="mx-auto" placeholder="Username" value={friendRequestUsername} onChange={e => setFriendRequestUsername(e.target.value)}/>
                <input className="m-2 btn btn-primary" type="submit" value="Send"/>
            </form>

            <div className="row m-2">
                <div className="friend-box border border-dark m-5">
                    <h5>Friends</h5>
                    {(friends != null && friends.friends != null) && friends.friends.map(function(f, i) {
                        // fix this
                        return <p><a href="javascript:void(0);" onClick={() => setRecpientID(f.user_id)}>@{f.username}</a> <a href="javascript:void(0);" onClick={removeFriend.bind(removeFriend, f.id)}>Remove</a></p>
                    })}
                    <h5>Incoming</h5>
                    {(friends != null && friends.incoming != null) && friends.incoming.map(function(f, i) {
                        return <p>@{f.sender_username} <a href="javascript:void(0);" onClick={acceptRequest.bind(acceptRequest, f.id)}>Accept</a> <a href="javascript:void(0);" onClick={declineRequest.bind(declineRequest, f.id)}>Decline</a></p>
                    })}
                    <h5>Pending</h5>
                    {(friends != null && friends.pending != null) && friends.pending.map(function(f, i) {
                        return <p>@{f.sender_username} <a href="javascript:void(0);" onClick={removeFriendRequest.bind(removeFriendRequest, f.id)}>Remove</a></p>
                    })}
                </div>

                <div className="chat-box justify-content-center border border-dark m-5">
                    <p>Hello</p>
                    <p>Hello how are you bro</p>
                </div>
            </div>

            <form className="text-center m-5" onSubmit={sendMessage}>
                <input className="m-4 mx-auto m-3" placeholder="Message" value={message} onChange={e => setMessage(e.target.value)}/>
                <input className="m-2 btn btn-primary" type="submit" value="Send"/>
            </form>
        </>
    )
}
