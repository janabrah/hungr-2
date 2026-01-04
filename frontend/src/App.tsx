import { useState } from 'react'
import { Home } from './pages/Home'
import { Upload } from './pages/Upload'
import { Browse } from './pages/Browse'

type Page = 'home' | 'upload' | 'browse'

function App() {
  const [page, setPage] = useState<Page>('home')

  switch (page) {
    case 'home':
      return <Home onNavigate={setPage} />
    case 'upload':
      return <Upload onNavigate={setPage} />
    case 'browse':
      return <Browse onNavigate={setPage} />
  }
}

export default App
