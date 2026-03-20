#!/usr/bin/env python
import psycopg2
import os
import sys
import time

user = os.environ.get('USER', 'postgres')
max_retries = 3

for attempt in range(max_retries):
    try:
        print(f"Attempt {attempt + 1} to create database...")
        conn = psycopg2.connect(f"dbname=postgres user={user}")
        conn.autocommit = True
        cur = conn.cursor()
        
        cur.execute("DROP DATABASE IF EXISTS krasnodar_tourism;")
        cur.execute("CREATE DATABASE krasnodar_tourism;")
        cur.close()
        conn.close()
        
        # Wait a moment for the database to be created
        time.sleep(1)
        
        conn = psycopg2.connect(f"dbname=krasnodar_tourism user={user}")
        conn.autocommit = True
        cur = conn.cursor()
        
        try:
            cur.execute("CREATE EXTENSION IF NOT EXISTS postgis;")
            print("✓ PostGIS extension enabled")
        except Exception as postgis_err:
            print(f"⚠ PostGIS not available: {postgis_err}")
            print("  Proceeding without PostGIS for now...")
        
        cur.close()
        conn.close()
        
        print("✓ Database created successfully")
        sys.exit(0)
    except Exception as e:
        print(f"✗ Error (attempt {attempt + 1}): {e}")
        if attempt < max_retries - 1:
            time.sleep(2)
            print("  Retrying...")
        else:
            print("✗ Failed to create database after all retries")
            sys.exit(1)
