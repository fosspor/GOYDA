#!/usr/bin/env python3
import socket
import time

def check_port(port, name):
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        result = sock.connect_ex(('localhost', port))
        sock.close()
        if result == 0:
            print(f"✅ {name} на порте {port} - ДОСТУПЕН")
            return True
        else:
            print(f"❌ {name} на порте {port} - НЕ ДОСТУПЕН")
            return False
    except Exception as e:
        print(f"⚠️  {name} - ошибка проверки: {e}")
        return False

print("Проверка доступности серверов...")
print("=" * 50)

time.sleep(10)

backend_ok = check_port(8000, "Django Backend")
frontend_ok = check_port(3000, "React Frontend")

print("\n" + "=" * 50)
if backend_ok and frontend_ok:
    print("\n🎉 ОБА СЕРВЕРА РАБОТАЮТ КОРРЕКТНО!\n")
    print("📋 Доступные URL:")
    print("   🔧 Django API:  http://localhost:8000")
    print("   👨‍💼 Admin Panel:  http://localhost:8000/admin")
    print("   💻 Frontend:    http://localhost:3000\n")
    print("Учетные данные для админ:")
    print("   Логин: admin")
    print("   Пароль: admin123456")
elif backend_ok:
    print("\n✅ Django Backend работает на порте 8000")
    print("⏳ React Frontend еще инициализируется (может занять до 30 секунд)...")
else:
    print("\n⚠️  Django Backend не отвечает!")
    print("❌ React Frontend не отвечает!")
