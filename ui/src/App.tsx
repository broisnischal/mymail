import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './store/auth'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import MailboxList from './pages/MailboxList'
import EmailList from './pages/EmailList'
import EmailView from './pages/EmailView'
import Layout from './components/Layout'

function App() {
  const { token } = useAuthStore()

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={!token ? <Login /> : <Navigate to="/" />} />
        <Route path="/register" element={!token ? <Register /> : <Navigate to="/" />} />
        <Route
          path="/"
          element={token ? <Layout /> : <Navigate to="/login" />}
        >
          <Route index element={<Dashboard />} />
          <Route path="mailboxes" element={<MailboxList />} />
          <Route path="emails" element={<EmailList />} />
          <Route path="emails/:id" element={<EmailView />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
