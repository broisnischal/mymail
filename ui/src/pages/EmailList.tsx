import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import api from '../api/client'

export default function EmailList() {
  const [emails, setEmails] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [searchParams] = useSearchParams()
  const mailboxId = searchParams.get('mailboxId')

  useEffect(() => {
    fetchEmails()
  }, [mailboxId])

  const fetchEmails = async () => {
    try {
      const params = mailboxId ? { mailboxId } : {}
      const { data } = await api.get('/emails', { params })
      setEmails(data.emails)
    } catch (err) {
      console.error('Failed to fetch emails:', err)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="p-6">Loading...</div>
  }

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">Emails</h1>
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="divide-y divide-gray-200">
          {emails.map((email) => (
            <Link
              key={email.id}
              to={`/emails/${email.id}`}
              className="block p-6 hover:bg-gray-50 transition-colors"
            >
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <p className="font-semibold text-gray-900">
                    {email.subject || '(No subject)'}
                  </p>
                  <p className="text-sm text-gray-600 mt-1">From: {email.from}</p>
                  <p className="text-sm text-gray-600">
                    To: {Array.isArray(email.to) ? email.to.join(', ') : email.to}
                  </p>
                </div>
                <span className="text-sm text-gray-500 ml-4">
                  {new Date(email.receivedAt).toLocaleString()}
                </span>
              </div>
            </Link>
          ))}
        </div>
        {emails.length === 0 && (
          <div className="p-6 text-center text-gray-500">No emails found</div>
        )}
      </div>
    </div>
  )
}
