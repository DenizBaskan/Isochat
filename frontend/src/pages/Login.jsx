import { useState } from 'react'
import axios from 'axios'
import { domain, hcaptcha_sitekey } from '../Globals'
import HCaptcha from '@hcaptcha/react-hcaptcha'
import Navbar from './Navbar'
import React from 'react'

export default function Login() {
    document.title = "Login"

    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [token, setToken] = useState("")

    const [err, setErr] = useState("")

    var captcha = React.createRef()

    const isEmail = (e) => {
        var re = /\S+@\S+\.\S+/
        return re.test(e)
    }

    const submit = (event) => {
        event.preventDefault();

        if (!isEmail(email)) {
            return setErr("Email is invalid")
        }

        if (password.length < 8) {
            return setErr("Password must be at least eight characters")
        }

        axios.post(domain + "/login", {
            "email": email,
            "captcha_key": token,
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
            <div className="custom-rounded m-5 p-3 card-login border-radius-5 text-center mx-auto">
                <h2 className="p-3">Login</h2>

                <p className="text-danger">{err}</p>

                <form onSubmit={submit}>
                    <input className="input mx-auto form-control m-3" placeholder="Email" onChange={e => setEmail(e.target.value)} />
                    <input className="input mx-auto form-control m-3" placeholder="Password" type="password" onChange={e => setPassword(e.target.value)}/>

                    <HCaptcha
                        sitekey={hcaptcha_sitekey}
                        onVerify={setToken}
                        onError={ () => { window.location.reload() } }
                        onExpire={ () => { window.location.reload() } }
                        ref={captcha}
                    />

                    <input className="m-2 btn btn-primary" type="submit"/>
                </form>

                <p className="m-2">Need to create an account? <a href="/signup">Signup</a></p>
            </div>
        </>
    )
}
