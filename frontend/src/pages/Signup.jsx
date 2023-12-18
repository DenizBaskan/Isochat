import { useState } from 'react'
import useCountdown from '../hooks/Countdown'
import axios from 'axios'
import { domain, hcaptcha_sitekey } from '../Globals'
import HCaptcha from '@hcaptcha/react-hcaptcha'
import Navbar from './Navbar'
import React from 'react'

export default function Signup() {
    document.title = "Signup"

    const [email, setEmail] = useState("")
    const [code, setCode] = useState("")
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")
    const [token, setToken] = useState("")

    const { secondsLeft, start } = useCountdown()

    const [err, setErr] = useState("")

    var captcha = React.createRef()

    const isEmail = (e) => {
        var re = /\S+@\S+\.\S+/
        return re.test(e)
    }

    const sendCode = () => {
        if (!isEmail(email)) {
            return setErr("Email is invalid")
        }
        
        axios.post(domain + "/signup/email", {
            "email": email
        }, {withCredentials : true}).then(() => {
            start(15)
        }).catch((e) => {
            if (e.response) {
                if (e.response.status == 400) {
                    return setErr(e.response.data["reason"])
                }
            }

            return setErr("An unexpected error occurred")
        })

        setErr("")
    }

    const submit = (event) => {
        event.preventDefault();

        if (!isEmail(email)) {
            return setErr("Email is invalid")
        }

        if (code.length != 6) {
            return setErr("Code must be six digits")
        }

        if (username.length == 0) {
            return setErr("Username must be at least one character")
        }

        if (password.length < 8) {
            return setErr("Password must be at least eight characters")
        }

        if (password != confirmPassword) {
            return setErr("Passwords do not match")
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
            captcha.current.resetCaptcha()
            
            if (e.response) {
                if (e.response.status == 400) {
                    return setErr(e.response.data["reason"])
                }
            }

            return setErr("An unexpected error occurred")
        })

        setErr("")
    }

    return (
        <>
            <Navbar />
            <div className="custom-rounded m-5 p-3 card-signup border-radius-5 text-center mx-auto">
                <h3 className="p-1 m-3">Signup</h3>

                <p className="text-danger">{err}</p>

                <form onSubmit={submit}>
                    <input className="input mx-auto form-control" placeholder="Email" onChange={e => setEmail(e.target.value)} />
                    <button className="btn btn-primary btn m-3" type="button" disabled={secondsLeft > 0} onClick={sendCode}>Send code {secondsLeft > 0 ? `(${secondsLeft})` : ""}</button>

                    <input className="input mx-auto form-control" placeholder="Code" onChange={e => setCode(e.target.value)}/>
        
                    <input placeholder="Username" className="input mx-auto form-control m-3" onChange={e => setUsername(e.target.value)}/>
                    <input placeholder="Password" className="input mx-auto form-control m-3" type="password" onChange={e => setPassword(e.target.value)}/>

                    <input className="input mx-auto form-control m-3" placeholder="Confirm password" type="password" onChange={e => setConfirmPassword(e.target.value)}/>

                    <HCaptcha
                        sitekey={hcaptcha_sitekey}
                        onVerify={setToken}
                        onError={ () => { window.location.reload() } }
                        onExpire={ () => { window.location.reload() } }
                        ref={captcha}
                    />

                    <input className="btn btn-primary m-4" type="submit" value="Submit"/>
                </form>

                <p>Already have an account? <a href="/login">Login</a></p>
            </div>
        </>
    )
}
