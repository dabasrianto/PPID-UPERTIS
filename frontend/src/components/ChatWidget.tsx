import React, { useState, useRef, useEffect } from 'react';
import { MessageSquare, X, RotateCcw, PhoneCall } from 'lucide-react';

export interface FAQItem {
  id?: number;
  question: string;
  answer: string;
}

const DEFAULT_FAQ_DATA: FAQItem[] = [
  {
    question: "Bagaimana cara mengajukan permohonan informasi?",
    answer: "Anda dapat mengajukan secara online melalui menu 'Permohonan Informasi' di portal ini, atau datang langsung ke Desk Layanan PPID UPERTIS di Gedung Rektorat Lantai 1. Siapkan identitas diri seperti KTP (perorangan) atau Akta Pendirian (organisasi)."
  },
  {
    question: "Berapa lama waktu proses permohonan informasi?",
    answer: "Sesuai dengan UU KIP No. 14 Tahun 2008, PPID akan memberikan jawaban/tanggapan dalam waktu 10 hari kerja sejak permohonan terdaftar, dan dapat diperpanjang paling lambat 7 hari kerja berikutnya."
  },
  {
    question: "Bagaimana jika permohonan informasi saya ditolak?",
    answer: "Jika permohonan ditolak atau tanggapan kurang memuaskan, Anda dapat mengajukan Keberatan Informasi secara online melalui menu 'Keberatan Informasi' di portal ini dalam waktu maksimal 30 hari kerja."
  },
  {
    question: "Apakah layanan informasi ini dipungut biaya?",
    answer: "Layanan permohonan informasi publik di PPID UPERTIS sepenuhnya GRATIS. Jika ada biaya penggandaan berkas fisik atau pengiriman dokumen, biaya tersebut ditanggung oleh pemohon."
  },
  {
    question: "Kapan jadwal operasional Desk Layanan PPID?",
    answer: "Desk fisik PPID buka Senin s/d Kamis (08:00 - 16:00 WIB) dan Jumat (08:00 - 16:30 WIB). Untuk layanan online melalui portal ini aktif 24/7."
  }
];

interface ChatMessage {
  sender: 'bot' | 'user';
  text: string;
  isFAQList?: boolean;
}

interface ChatWidgetProps {
  rektoratPhone?: string;
  faqs?: FAQItem[];
}

export default function ChatWidget({ rektoratPhone, faqs }: ChatWidgetProps) {
  const activeFaqs = faqs && faqs.length > 0 ? faqs : DEFAULT_FAQ_DATA;
  const [isOpen, setIsOpen] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([
    {
      sender: 'bot',
      text: 'Halo! Selamat datang di Layanan Chatbot FAQ PPID UPERTIS. Ada yang bisa kami bantu hari ini? Silakan klik salah satu pertanyaan yang sering ditanyakan di bawah ini:',
      isFAQList: true
    }
  ]);

  const messageEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (messageEndRef.current) {
      messageEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages, isOpen]);


  const handleFAQClick = (faq: FAQItem) => {
    // Add user message
    const userMsg: ChatMessage = {
      sender: 'user',
      text: faq.question
    };

    // Add bot answer message
    const botMsg: ChatMessage = {
      sender: 'bot',
      text: faq.answer
    };

    setMessages(prev => [...prev, userMsg, botMsg]);
  };

  const handleReset = () => {
    setMessages([
      {
        sender: 'bot',
        text: 'Halo! Selamat datang di Layanan Chatbot FAQ PPID UPERTIS. Ada yang bisa kami bantu hari ini? Silakan klik salah satu pertanyaan yang sering ditanyakan di bawah ini:',
        isFAQList: true
      }
    ]);
  };

  const formattedPhone = rektoratPhone
    ? (rektoratPhone as string).replace(/[^0-9]/g, '').replace(/^0/, '62')
    : '';

  return (
    <div className="fixed bottom-20 lg:bottom-5 right-5 z-50 flex flex-col items-end">
      {/* Chat Window */}
      {isOpen && (
        <div className="mb-4 w-[350px] sm:w-[380px] h-[500px] bg-white rounded-[2.5rem] border border-slate-200 shadow-2xl flex flex-col overflow-hidden animate-in slide-in-from-bottom duration-300">
          {/* Header */}
          <div className="bg-gradient-to-r from-[#002147] to-[#0b335c] p-5 text-white flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="relative">
                <div className="w-10 h-10 rounded-full bg-amber-400/20 border border-amber-400 flex items-center justify-center font-bold text-amber-400">
                  PPID
                </div>
                <span className="absolute bottom-0 right-0 w-2.5 h-2.5 bg-green-500 rounded-full border-2 border-[#002147]"></span>
              </div>
              <div className="text-left">
                <h4 className="text-xs font-bold leading-tight">Asisten FAQ PPID</h4>
                <p className="text-[9px] text-slate-300 font-medium">Aktif untuk membantu Anda</p>
              </div>
            </div>
            <button
              onClick={() => setIsOpen(false)}
              className="w-8 h-8 rounded-full bg-white/10 hover:bg-white/20 flex items-center justify-center transition-colors cursor-pointer"
            >
              <X className="h-4 w-4 text-white" />
            </button>
          </div>

          {/* Chat Messages Area */}
          <div className="flex-1 p-5 overflow-y-auto space-y-4 bg-slate-50 text-xs">
            {messages.map((msg, index) => (
              <div
                key={index}
                className={`flex ${msg.sender === 'user' ? 'justify-end' : 'justify-start'}`}
              >
                <div
                  className={`max-w-[85%] p-4 rounded-3xl text-left leading-relaxed ${
                    msg.sender === 'user'
                      ? 'bg-amber-500 text-white font-semibold rounded-tr-none'
                      : 'bg-white border border-slate-200 text-slate-700 font-medium rounded-tl-none shadow-sm'
                  }`}
                >
                  <p>{msg.text}</p>

                  {/* Render FAQ questions only on the bot greeting or inside the message list if requested */}
                  {msg.isFAQList && (
                    <div className="mt-4 space-y-2 border-t border-slate-100 pt-3">
                      {activeFaqs.map((faq, idx) => (
                        <button
                          key={idx}
                          onClick={() => handleFAQClick(faq)}
                          className="w-full text-left p-2.5 rounded-xl border border-slate-100 hover:border-amber-400 hover:bg-amber-50/50 text-[11px] text-slate-700 font-bold transition-all cursor-pointer block bg-slate-50/50"
                        >
                          ❓ {faq.question}
                        </button>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            ))}
            <div ref={messageEndRef} />
          </div>

          {/* Action Bar (Reset / Contact WA) */}
          <div className="p-4 bg-white border-t border-slate-100 flex gap-2">
            <button
              onClick={handleReset}
              className="flex-1 flex items-center justify-center gap-1.5 py-2.5 px-4 bg-slate-150 hover:bg-slate-200 text-[#002147] font-bold rounded-2xl transition-all cursor-pointer text-[11px] border border-slate-200"
            >
              <RotateCcw className="h-3.5 w-3.5" />
              Mulai Ulang
            </button>
            {formattedPhone && (
              <a
                href={`https://wa.me/${formattedPhone}`}
                target="_blank"
                rel="noopener noreferrer"
                className="flex-1 flex items-center justify-center gap-1.5 py-2.5 px-4 bg-[#25D366] hover:bg-[#20ba5a] text-white font-bold rounded-2xl transition-all cursor-pointer text-[11px] shadow-sm hover:shadow-green-500/20"
              >
                <PhoneCall className="h-3.5 w-3.5" />
                CS WhatsApp
              </a>
            )}
          </div>
        </div>
      )}

      {/* Floating Toggle Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-[52px] h-[52px] rounded-full bg-amber-500 text-white shadow-lg shadow-amber-500/30 flex items-center justify-center hover:scale-110 hover:bg-amber-600 hover:shadow-amber-500/50 transition-all duration-300 cursor-pointer active:scale-95 z-50 border-0"
      >
        {isOpen ? (
          <X className="h-6 w-6" />
        ) : (
          <MessageSquare className="h-6 w-6" />
        )}
      </button>
    </div>
  );
}
