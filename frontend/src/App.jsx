import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { Navigate } from 'react-router-dom'

import './index.css'

import Chat from './pages/Chat'
import Signup from './pages/Signup'
import Login from './pages/Login'
import NoPage from './pages/NoPage'
import About from './pages/About'
import Contact from './pages/Contact'

export default function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Navigate to="/chat" />} />
                <Route path="/chat" element={<Chat />} />
                <Route path="/signup" element={<Signup />} />
                <Route path="/login" element={<Login />} />
                <Route path="/about" element={<About />} />
                <Route path="/contact" element={<Contact />} />
                <Route path="*" element={<NoPage />} />
            </Routes>
        </BrowserRouter>
    )
}
