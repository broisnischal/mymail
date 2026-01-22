import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import api from '../api/client'

export default function MailboxList() {
  const [mailboxes, setMailboxes] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [newAddress, setNewAddress] = useState('')
  const [isAlias, setIsAlias] = useState(false)

  useEffect(() => {
    fetchMailboxes()
  }, [])

  const fetchMailboxes = async () => {
    try {
      const { data } = await api.get('/mailboxes')
      setMailboxes(data.mailboxes)
    } catch (err) {
      console.error('Failed to fetch mailboxes:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.post('/mailboxes', { address: newAddress, isAlias })
      setNewAddress('')
      setIsAlias(false)
      fetchMailboxes()
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to create mailbox')
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this mailbox?')) return
    try {
      await api.delete(`/mailboxes/${id}`)
      fetchMailboxes()
    } catch (err) {
      alert('Failed to delete mailbox')
    }
  }

  if (loading) {
    return <div className="p-6">Loading...</div>
  }

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Mailboxes</h1>
      </div>

      <form onSubmit={handleCreate} className="bg-white p-6 rounded-lg shadow mb-6">
        <h2 className="text-xl font-semibold mb-4">Create Mailbox</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email Address
            </label>
            <input
              type="email"
              value={newAddress}
              onChange={(e) => setNewAddress(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
              placeholder="user@mymail.com"
              required
            />
          </div>
          <div className="flex items-center">
            <input
              type="checkbox"
              id="isAlias"
              checked={isAlias}
              onChange={(e) => setIsAlias(e.target.checked)}
              className="mr-2"
            />
            <label htmlFor="isAlias" className="text-sm text-gray-700">
              Is Alias
            </label>
          </div>
          <button
            type="submit"
            className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
          >
            Create
          </button>
        </div>
      </form>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Address
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Type
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Created
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {mailboxes.map((mailbox) => (
              <tr key={mailbox.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {mailbox.address}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {mailbox.isTemp ? 'Temp' : mailbox.isAlias ? 'Alias' : 'Primary'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(mailbox.createdAt).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <Link
                    to={`/emails?mailboxId=${mailbox.id}`}
                    className="text-indigo-600 hover:text-indigo-900 mr-4"
                  >
                    View Emails
                  </Link>
                  <button
                    onClick={() => handleDelete(mailbox.id)}
                    className="text-red-600 hover:text-red-900"
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
