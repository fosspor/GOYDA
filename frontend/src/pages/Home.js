import React from 'react';

export default function Home() {
  return (
    <main className="min-h-screen">
      {/* Hero Section */}
      <section className="bg-gradient-to-r from-blue-600 to-purple-600 text-white py-20">
        <div className="container mx-auto px-4 text-center">
          <h1 className="text-5xl font-bold mb-4">GOYDA</h1>
          <p className="text-xl mb-8">Открой Краснодарский край по-новому</p>
          <p className="text-lg max-w-2xl mx-auto mb-8">
            Персональные маршруты, местные сокровища и незабываемые впечатления
          </p>
          <button className="bg-white text-blue-600 px-8 py-3 rounded-lg font-bold hover:bg-gray-100 transition">
            Начать путешествие
          </button>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-white">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-12">Наши возможности</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {[
              {
                title: 'Персональные маршруты',
                description: 'AI подберет идеальный маршрут под ваши интересы'
              },
              {
                title: 'Локальные сокровища',
                description: 'Откройте заведения, которые не найдете в туристических гайдах'
              },
              {
                title: 'Виртуальные туры',
                description: 'Посетите локации в 360° формате перед поездкой'
              },
              {
                title: 'Сезонные рекомендации',
                description: 'Получайте советы в зависимости от времени года'
              },
              {
                title: 'Поддержка бизнеса',
                description: 'Маленькие ферм и винодельни получают больше гостей'
              },
              {
                title: 'Логистика включена',
                description: 'Помощь с транспортом, ночлегом и питанием'
              },
            ].map((feature, idx) => (
              <div key={idx} className="p-6 border rounded-lg hover:shadow-lg transition">
                <h3 className="text-xl font-bold mb-2">{feature.title}</h3>
                <p className="text-gray-600">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>
    </main>
  );
}
