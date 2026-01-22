import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import api from '../api/client'

export default function EmailView() {
  const { id } = useParams()
  const [email, setEmail] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (id) {
      fetchEmail()
    }
  }, [id])

  const fetchEmail = async () => {
    try {
      const { data } = await api.get(`/emails/${id}`)
      setEmail(data.email)
    } catch (err) {
      console.error('Failed to fetch email:', err)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="p-6">Loading...</div>
  }

  if (!email) {
    return <div className="p-6">Email not found</div>
  }

  return (
    <div className="p-6">
      <Link to="/emails" className="text-indigo-600 hover:text-indigo-900 mb-4 inline-block">
        ← Back to emails
      </Link>
      <div className="bg-white rounded-lg shadow p-6">
        <div className="border-b pb-4 mb-4">
          <h1 className="text-2xl font-bold mb-2">{email.subject || '(No subject)'}</h1>
          <div className="text-sm text-gray-600 space-y-1">
            <p><strong>From:</strong> {email.from}</p>
            <p><strong>To:</strong> {Array.isArray(email.to) ? email.to.join(', ') : email.to}</p>
            {email.cc && email.cc.length > 0 && (
              <p><strong>CC:</strong> {Array.isArray(email.cc) ? email.cc.join(', ') : email.cc}</p>
            )}
            <p><strong>Date:</strong> {new Date(email.receivedAt).toLocaleString()}</p>
          </div>
        </div>
        <div className="prose max-w-none">
          {email.htmlBody ? (
            <div dangerouslySetInnerHTML={{ __html: email.htmlBody }} />
          ) : (
            <pre className="whitespace-pre-wrap font-sans">{email.textBody || 'No content'}</pre>
          )}
        </div>
      </div>
    </div>
  )
}
