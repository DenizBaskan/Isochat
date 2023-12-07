import { useState } from 'react'
import useCountdown from '../hooks/Countdown'
import axios from 'axios'
import { domain, hcaptcha_sitekey } from '../Globals'
import HCaptcha from '@hcaptcha/react-hcaptcha'
import { Navigate } from 'react-router-dom'

export default function Signup() {
    document.title = "Signup"

    const [email, setEmail] = useState("")
    const [code, setCode] = useState("")
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")
    const [token, setToken] = useState("")

    const { secondsLeft, start } = useCountdown()

    const isEmail = (e) => {
        var re = /\S+@\S+\.\S+/
        return re.test(e)
    }

    const sendCode = () => {
        if (!isEmail(email)) {
            return alert("Email is invalid")
        }
        
        axios.post(domain + "/signup/email", {
            "email": email
        }, {withCredentials : true}).then(() => {
            start(15)
        }).catch((e) => {
            if (e.response) {
                if (e.response.status == 400) {
                    return alert(e.response.data["reason"])
                }
            }

            alert("An unexpected error occurred")
        })
    }

    const submit = (event) => {
        event.preventDefault();

        if (!isEmail(email)) {
            return alert("Email is invalid")
        }

        if (code.length != 6) {
            return alert("Code must be six digits")
        }

        if (code.length != 6) {
            return alert("Code must be six digits")
        }

        if (username.length == 0) {
            return alert("Username must be at least one character")
        }

        if (password.length < 8) {
            return alert("Password must be at least eight characters")
        }

        if (password != confirmPassword) {
            return alert("Passwords do not match")
        }

        axios.post(domain + "/signup", {
            "email": email,
            "captcha_key": token,
            "code": code,
            "username": username,
            "password": password
        }, {withCredentials : true}).then(() => {
            window.location.href = "/chat"
        }).catch((e) => {
            if (e.response) {
                if (e.response.status == 400) {
                    return alert(e.response.data["reason"])
                }
            }

            return alert("An unexpected error occurred")
        })
    }

    return (
        <>
            <h1>Signup</h1>

            <form onSubmit={submit}>
                <label>Email</label>
                <input onChange={e => setEmail(e.target.value)} />

                <button type="button" disabled={secondsLeft > 0} onClick={sendCode}>Send code {secondsLeft > 0 ? `(${secondsLeft})` : ""}</button>

                <br />
                <br />

                <label>Code</label>
                <input onChange={e => setCode(e.target.value)}/>

                <br />

                <label>Username</label>
                <input onChange={e => setUsername(e.target.value)}/>

                <br />

                <label>Password</label>
                <input type="password" onChange={e => setPassword(e.target.value)}/>

                <br />

                <label>Confirm password</label>
                <input type="password" onChange={e => setConfirmPassword(e.target.value)}/>

                <br />
                <br />

                <HCaptcha
                    sitekey={hcaptcha_sitekey}
                    onVerify={setToken}
                    onError={ () => { window.location.reload() } }
                    onExpire={ () => { window.location.reload() } }
                />

                <input type="submit"/>
            </form>

            <p>Already have an account? <a href="/login">Login</a></p>
        </>
    )
}
