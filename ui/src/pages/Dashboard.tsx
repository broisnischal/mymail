import { useEffect, useState } from 'react'
import api from '../api/client'

export default function Dashboard() {
  const [stats, setStats] = useState({
    mailboxes: 0,
    emails: 0,
    recentEmails: [] as any[],
  })

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [mailboxesRes, emailsRes] = await Promise.all([
          api.get('/mailboxes'),
          api.get('/emails?limit=10'),
        ])
        setStats({
          mailboxes: mailboxesRes.data.mailboxes.length,
          emails: emailsRes.data.emails.length,
          recentEmails: emailsRes.data.emails,
        })
      } catch (err) {
        console.error('Failed to fetch stats:', err)
      }
    }
    fetchStats()
  }, [])

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold mb-2">Mailboxes</h2>
          <p className="text-3xl font-bold text-indigo-600">{stats.mailboxes}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h2 className="text-xl font-semibold mb-2">Total Emails</h2>
          <p className="text-3xl font-bold text-indigo-600">{stats.emails}</p>
        </div>
      </div>
      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          <h2 className="text-xl font-semibold mb-4">Recent Emails</h2>
          <div className="space-y-4">
            {stats.recentEmails.map((email) => (
              <div key={email.id} className="border-b pb-4">
                <div className="flex justify-between items-start">
                  <div>
                    <p className="font-semibold">{email.subject || '(No subject)'}</p>
                    <p className="text-sm text-gray-600">From: {email.from}</p>
                    <p className="text-sm text-gray-600">To: {email.to.join(', ')}</p>
                  </div>
                  <span className="text-sm text-gray-500">
                    {new Date(email.receivedAt).toLocaleDateString()}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
