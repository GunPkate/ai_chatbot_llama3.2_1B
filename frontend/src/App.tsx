import { useState, useRef, useEffect, type KeyboardEvent } from "react";

const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080/chat";

type Message = {
  role: "user" | "assistant";
  content: string;
};

export default function App() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);

  // Scroll to the newest message whenever the list changes
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // Append text to the last (assistant) message
  function appendToLast(token: string) {
    setMessages((prev) => {
      if (prev.length === 0) return prev;
      const copy = [...prev];
      const last = copy[copy.length - 1];
      copy[copy.length - 1] = { ...last, content: last.content + token };
      return copy;
    });
  }

  async function sendMessage() {
    const text = input.trim();
    if (!text || loading) return;

    setInput("");
    setLoading(true);
    // Add the user's message plus an empty assistant message to fill in
    setMessages((prev) => [
      ...prev,
      { role: "user", content: text },
      { role: "assistant", content: "" },
    ]);

    try {
      const res = await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ message: text }),
      });
      if (!res.ok || !res.body) throw new Error(`HTTP ${res.status}`);

      // Read the response body as a stream of bytes
      const reader = res.body.getReader();
      const decoder = new TextDecoder();
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });

        // Events are separated by a blank line. The last piece may be
        // incomplete, so keep it in the buffer for the next round.
        const events = buffer.split("\n\n");
        buffer = events.pop() ?? "";

        for (const event of events) {
          if (!event.startsWith("data: ")) continue;
          appendToLast(JSON.parse(event.slice(6)) as string);
        }
      }
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err);
      appendToLast(`\n[Error: ${msg}]`);
    } finally {
      setLoading(false);
    }
  }

  function handleKeyDown(e: KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  }

  return (
    <div className="mx-auto flex h-screen max-w-2xl flex-col p-4">
      <h1 className="mb-2 text-xl font-semibold">Local AI Chatbot</h1>

      <div className="flex flex-1 flex-col gap-2.5 overflow-y-auto py-2">
        {messages.length === 0 && (
          <p className="mt-10 text-center text-gray-500">
            Ask something to get started.
          </p>
        )}
        {messages.map((m, i) => (
          <div
            key={i}
            className={
              "max-w-[80%] whitespace-pre-wrap rounded-xl px-3.5 py-2.5 leading-snug " +
              (m.role === "user"
                ? "self-end bg-blue-600 text-white"
                : "self-start bg-gray-200 text-gray-900")
            }
          >
            {m.content || (loading ? "…" : "")}
          </div>
        ))}
        <div ref={bottomRef} />
      </div>

      <div className="flex gap-2">
        <textarea
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Type a message (Enter to send)"
          rows={2}
          className="flex-1 resize-none rounded-lg border border-gray-300 p-2.5 focus:border-blue-600 focus:outline-none"
        />
        <button
          onClick={sendMessage}
          disabled={loading || !input.trim()}
          className="rounded-lg bg-blue-600 px-5 text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
        >
          Send
        </button>
      </div>
    </div>
  );
}