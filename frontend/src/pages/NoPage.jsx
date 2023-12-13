import Navbar from './Navbar'

export default function NoPage() {
    document.title = "404"
    
    return (
        <>
            <Navbar />

            <h3 className="text-center m-3">Chat</h3>
        </>
    )
}
