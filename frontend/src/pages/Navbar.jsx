import axios from 'axios'
import { domain } from '../Globals'
import { useState, useEffect } from 'react'

export default function Navbar() {
    const AuthStatus = {
        Loading: 1,
        Authorized: 2,
        Unauthorized: 3
    }
    
    const [authed, setAuthed] = useState(AuthStatus.Loading)

    useEffect(() => {
        axios.post(domain + "/authed", undefined, {withCredentials : true})
        .then(() => {
            setAuthed(AuthStatus.Authorized)

            if (window.location.pathname == "/login" || window.location.pathname == "/signup") {
                window.location.href = "/chat"
            }
        })
        .catch(() => {
            setAuthed(AuthStatus.Unauthorized)

            if (window.location.pathname == "/chat") {
                window.location.href = "/login"
            }
        })
    }, []);

    const logout = () => {
        axios.post(domain + "/logout", undefined, {withCredentials : true})
        .then(() => {
            window.location.href = "/login"
        })
        .catch(() => {
            if (e.response) {
                if (e.response.status == 400) {
                    return alert(e.response.data["reason"])
                }
            }

            alert("An unexpected error occurred")
        })
    }

    return (
        <nav className="navbar navbar-col navbar navbar-expand p-3">
            <div className="container-fluid">
                <a className="navbar-a navbar-brand" href="/about">Isochat</a>
                <div>
                    <ul className="navbar-nav ms-auto ">
                        <li className="nav-item">
                            <a className="navbar-a nav-link mx-2" href="/about">About</a>
                        </li>
                        <li className="nav-item">
                            <a className="navbar-a nav-link mx-2" href="/contact">Contact</a>
                        </li>
                        {authed == AuthStatus.Unauthorized && <li className="nav-item">
                            <a className="navbar-a nav-link mx-2" href="/login">Login</a>
                        </li>}
                        {authed == AuthStatus.Unauthorized && <li className="nav-item dropdown">
                            <a className="navbar-a nav-link mx-2" href="/signup">Signup</a>
                        </li>}
                        {authed == AuthStatus.Authorized && <li className="nav-item dropdown">
                            <a className="navbar-a nav-link mx-2" href="/chat">Chat</a>
                        </li>}
                        {authed== AuthStatus.Authorized && <li className="nav-item dropdown">
                            <a className="navbar-a nav-link mx-2 pointer" onClick={logout}>Logout</a>
                        </li>}
                    </ul>
                </div>
            </div>
        </nav>
    )
}
