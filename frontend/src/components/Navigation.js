import React from 'react';
import { Link } from 'react-router-dom';
import { FiHome, FiCompass, FiMap, FiUser, FiMenu, FiX } from 'react-icons/fi';
import './Navigation.css';

export default function Navigation() {
  const [isMenuOpen, setIsMenuOpen] = React.useState(false);

  return (
    <nav className="navbar bg-white shadow-sm sticky top-0 z-50">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link to="/" className="font-bold text-2xl text-blue-600">
            GOYDA
          </Link>

          {/* Desktop Menu */}
          <div className="hidden md:flex gap-8">
            <Link to="/" className="flex items-center gap-2 hover:text-blue-600 transition">
              <FiHome /> Главная
            </Link>
            <Link to="/discover" className="flex items-center gap-2 hover:text-blue-600 transition">
              <FiCompass /> Открыть локации
            </Link>
            <Link to="/suggested-routes" className="flex items-center gap-2 hover:text-blue-600 transition">
              <FiMap /> Мои маршруты
            </Link>
          </div>

          {/* Profile & Menu Toggle */}
          <div className="flex items-center gap-4">
            <Link to="/profile" className="hidden md:flex items-center gap-2 hover:text-blue-600 transition">
              <FiUser /> Профиль
            </Link>
            <button
              className="md:hidden text-2xl"
              onClick={() => setIsMenuOpen(!isMenuOpen)}
            >
              {isMenuOpen ? <FiX /> : <FiMenu />}
            </button>
          </div>
        </div>

        {/* Mobile Menu */}
        {isMenuOpen && (
          <div className="md:hidden pb-4 border-t">
            <Link to="/" className="block py-2 text-center hover:text-blue-600">
              Главная
            </Link>
            <Link to="/discover" className="block py-2 text-center hover:text-blue-600">
              Открыть локации
            </Link>
            <Link to="/suggested-routes" className="block py-2 text-center hover:text-blue-600">
              Мои маршруты
            </Link>
            <Link to="/profile" className="block py-2 text-center hover:text-blue-600">
              Профиль
            </Link>
          </div>
        )}
      </div>
    </nav>
  );
}
