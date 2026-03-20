# Frontend Installation & Development Guide

## Prerequisites

- Node.js 16+ (ideally 18 LTS)
- npm 7+ or yarn

## Installation

### 1. Install Node Dependencies

```bash
cd frontend
npm install
```

This will install:
- React 18 & React DOM
- React Router v6 for navigation
- Axios for API calls
- Zustand for state management
- Tailwind CSS for styling
- Leaflet for maps
- React Icons for UI icons
- And many more...

### 2. Environment Setup

```bash
# Create .env file (optional for development)
cp .env.example .env.local
```

Update with your backend URL:
```env
REACT_APP_API_URL=http://localhost:8000/api
```

## Development

### Start Dev Server

```bash
npm start
```

Opens at `http://localhost:3000` with hot reload enabled.

## Scripts

```bash
# Development server with hot reload
npm start

# Production build
npm run build

# Run tests
npm test

# Test coverage
npm test -- --coverage

# Eject config (irreversible!)
npm run eject
```

## Project Structure

```
src/
├── components/          # Reusable components
│   ├── Navigation.js    # Main navbar
│   ├── LocationCard.js  # Location display card
│   └── ...
├── pages/              # Full page components
│   ├── Home.js
│   ├── Discover.js
│   ├── SuggestedRoutes.js
│   └── ...
├── services/           # API clients
│   ├── api.js         # Axios setup & endpoints
│   └── auth.js        # Authentication helpers
├── store/             # Zustand state stores
│   └── index.js
├── hooks/             # Custom React hooks
├── App.js             # Root component
├── App.css
└── index.js           # Entry point
```

## Key Technologies

### Styling
- **Tailwind CSS** v3.3 for utility-first CSS
- Class-based styling with responsive design
- Mobile-first approach

### State Management
- **Zustand** for global state
- Simple, lightweight alternative to Redux

### HTTP Client
- **Axios** with interceptors
- JWT token handling
- Error handling middleware

### Maps
- **Leaflet** + **React-Leaflet**
- OpenStreetMap for base layer
- Markers for locations

### Routing
- **React Router v6**
- Nested routes support
- Dynamic route parameters

## Component Architecture

### Example: LocationCard Component

```javascript
export default function LocationCard({ location }) {
  return (
    <div className="border rounded-lg hover:shadow-lg transition">
      <img src={location.image} alt={location.name} />
      <h3 className="font-bold text-lg">{location.name}</h3>
      <p className="text-gray-600">{location.short_description}</p>
      <div className="flex justify-between items-center">
        <span className="text-yellow-500">★ {location.rating}</span>
        <Link to={`/location/${location.id}`}>Подробнее →</Link>
      </div>
    </div>
  );
}
```

## State Management Example

```javascript
// In your component
import { useLocationStore } from '../store';

export default function Discover() {
  const locations = useLocationStore(state => state.locations);
  const setLocations = useLocationStore(state => state.setLocations);
  const filters = useLocationStore(state => state.filters);
  
  // Use locations, setLocations, filters
}
```

## API Integration

### Making Requests

```javascript
import { locationsAPI } from '../services/api';

// In component
const [locations, setLocations] = useState([]);

useEffect(() => {
  locationsAPI.getAll({ category: 'wine' })
    .then(response => setLocations(response.data.results))
    .catch(error => console.error('Error:', error));
}, []);
```

### Authentication

```javascript
import { authAPI } from '../services/api';

const handleLogin = async (username, password) => {
  try {
    const response = await authAPI.login(username, password);
    localStorage.setItem('access_token', response.data.access);
    localStorage.setItem('refresh_token', response.data.refresh);
  } catch (error) {
    alert('Login failed');
  }
};
```

## Responsive Design

### Mobile-First Breakpoints (Tailwind)

```html
<!-- Default (mobile) -->
<div class="w-full">

<!-- Medium screens (tablet) -->
<div class="md:w-1/2">

<!-- Large screens (desktop) -->
<div class="lg:w-1/3">
```

## Browser Support

Targets:
- Chrome/Edge: latest 2 versions
- Firefox: latest 2 versions
- Safari: latest 2 versions
- Mobile browsers: iOS 12+, Android 5+

## Performance Tips

1. **Code Splitting**: React Router automatically splits by routes
2. **Image Optimization**: Use Tailwind `object-fit` for thumbnails
3. **Lazy Loading**: Use `React.lazy()` for heavy components
4. **Memoization**: Use `React.memo()` for expensive components

## Testing

### Example Test

```javascript
import { render, screen } from '@testing-library/react';
import Home from '../pages/Home';

test('renders hero section', () => {
  render(<Home />);
  expect(screen.getByText(/Краснодарский край/i)).toBeInTheDocument();
});
```

## Deployment

### Build Production Version

```bash
npm run build
```

Creates optimized build in `build/` directory.

### Deploy to Vercel (Recommended)

```bash
npm install -g vercel
vercel
```

### Deploy to Netlify

```bash
npm run build
# Then drag `build/` folder to Netlify
```

## Troubleshooting

### Port 3000 already in use
```bash
# Kill process on port 3000
lsof -i :3000 | grep LISTEN | awk '{print $2}' | xargs kill -9

# Or use different port
PORT=3001 npm start
```

### Module not found error
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### API requests failing
```bash
# Check if backend is running on port 8000
curl http://localhost:8000/api/locations/

# Check CORS settings in Django backend
```

## Useful Commands

```bash
# Check code quality
npm run lint

# Format code
npm run format

# Security audit
npm audit

# Update all packages
npm update
```

## Next.js Alternative

If you want to switch to Next.js (SSR/SSG):

```bash
# Create Next.js project
npx create-next-app@latest goyda --typescript

# Install additional packages
npm install zustand leaflet react-leaflet
```

---

**Happy coding! 🚀**
