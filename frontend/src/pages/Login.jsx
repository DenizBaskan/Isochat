import { useState } from 'react'
import axios from 'axios'
import { domain, hcaptcha_sitekey } from '../Globals'
import HCaptcha from '@hcaptcha/react-hcaptcha'

export default function Login() {
    document.title = "Login"

    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [token, setToken] = useState("")

    const isEmail = (e) => {
        var re = /\S+@\S+\.\S+/
        return re.test(e)
    }

    const submit = (event) => {
        event.preventDefault();

        if (!isEmail(email)) {
            return alert("Email is invalid")
        }

        if (password.length < 8) {
            return alert("Password must be at least eight characters")
        }

        axios.post(domain + "/login", {
            "email": email,
            "captcha_key": token,
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
            <h1>Login</h1>

            <form onSubmit={submit}>
                <label>Email</label>
                <input onChange={e => setEmail(e.target.value)} />

                <br />

                <label>Password</label>
                <input type="password" onChange={e => setPassword(e.target.value)}/>

                <HCaptcha
                    sitekey={hcaptcha_sitekey}
                    onVerify={setToken}
                    onError={ () => { window.location.reload() } }
                    onExpire={ () => { window.location.reload() } }
                />

                <input type="submit"/>
            </form>

            <p>Need to create an account? <a href="/signup">Signup</a></p>
        </>
    )
}
