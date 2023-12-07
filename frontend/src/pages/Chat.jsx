import axios from 'axios'
import { domain } from '../Globals'

export default function Chat() {
    document.title = "Chat"

    const logout = () => {
        axios.post(domain + "/logout", undefined, {withCredentials : true})
        .then(() => {
            window.location.href = "/login"
        })
        .catch(() => {
            alert("An unexpected error occurred")
        })
    }

    return (
        <>
            <h1>Chat</h1>

            <button onClick={logout}>Logout</button>
        </>
    )
}
