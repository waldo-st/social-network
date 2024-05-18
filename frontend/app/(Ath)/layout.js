"use client"
import { useRouter } from 'next/navigation';
import "../globals.css";


export default function RootLayout({ children }) {
  useRouter();
  return (
    <html lang="en">
      <body>
        <div>
          {children}
        </div>
      </body>
    </html>
  );
}
