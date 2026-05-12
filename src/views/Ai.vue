<template>
  <div class="chat-container">
    <!-- Header -->
    <div class="chat-header">
      <div class="header-content">
        <i class="pi pi-comments" style="font-size: 1.5rem; margin-right: 10px;"></i>
        <h1>AI Assistant</h1>
      </div>
      <Button 
        icon="pi pi-plus" 
        label="New Chat" 
        class="p-button-outlined" 
        @click="startNewChat"
      />
    </div>

    <!-- Messages Area -->
    <div class="messages-area" ref="messagesContainer">
      <div v-if="messages.length === 0" class="welcome-screen">
        <i class="pi pi-android" style="font-size: 4rem; color: var(--primary-color); margin-bottom: 1rem;"></i>
        <h2>How can I help you today?</h2>
        <p>Ask me anything! I'm here to assist with your questions.</p>
        
        <!-- Suggested prompts -->
        <div class="suggestions">
          <Chip label="Explain a concept" class="suggestion-chip" @click="setPrompt('Can you explain how machine learning works?')" />
          <Chip label="Write code" class="suggestion-chip" @click="setPrompt('Write a Python function to sort a list')" />
          <Chip label="Brainstorm ideas" class="suggestion-chip" @click="setPrompt('Give me 5 creative marketing ideas')" />
          <Chip label="Summarize text" class="suggestion-chip" @click="setPrompt('Summarize the key points of photosynthesis')" />
        </div>
      </div>

      <!-- Message List -->
      <div v-for="(message, index) in messages" :key="index" :class="['message', message.type]">
        <div class="message-avatar">
          <i :class="message.type === 'user' ? 'pi pi-user' : 'pi pi-android'"></i>
        </div>
        <div class="message-content">
          <div class="message-text" v-html="formatMessage(message.text)"></div>
          <div class="message-time">{{ message.timestamp }}</div>
        </div>
      </div>

      <!-- Typing Indicator -->
      <div v-if="isTyping" class="message bot">
        <div class="message-avatar">
          <i class="pi pi-android"></i>
        </div>
        <div class="message-content">
          <div class="typing-indicator">
            <span></span>
            <span></span>
            <span></span>
          </div>
        </div>
      </div>
    </div>

    <!-- Input Area -->
    <div class="input-area">
      <div class="input-container">
        <Textarea 
          v-model="currentPrompt" 
          placeholder="Send a message..." 
          :autoResize="true" 
          rows="1"
          class="chat-input"
          @keydown.enter.prevent="sendMessage"
          @keydown.enter.shift="addNewLine"
        />
        <Button 
          icon="pi pi-send" 
          class="send-button" 
          :loading="isTyping"
          @click="sendMessage"
          :disabled="!currentPrompt.trim() || isTyping"
        />
      </div>
      <div class="input-footer">
        <small>Press Enter to send, Shift + Enter for new line</small>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import Button from 'primevue/button'
import Textarea from 'primevue/textarea'
import Chip from 'primevue/chip'

// State
const messages = ref([])
const currentPrompt = ref('')
const isTyping = ref(false)
const messagesContainer = ref(null)

// Mock API call - Replace with your actual API endpoint
const sendToAI = async (prompt) => {
  // Simulate API delay
  await new Promise(resolve => setTimeout(resolve, 1500))
  
  // Mock response - Replace with actual API call
  const mockResponses = [
    `I understand you're asking about "${prompt}". That's a great question! Here's what I can tell you about it...`,
    `Based on your query "${prompt}", I think the best approach would be to consider the following points...`,
    `Thanks for your question about "${prompt}". Let me provide you with a comprehensive answer...`,
    `I've analyzed your prompt "${prompt}". Here's my response with some useful information...`
  ]
  
  return mockResponses[Math.floor(Math.random() * mockResponses.length)]
}

// Send message function
const sendMessage = async () => {
  if (!currentPrompt.value.trim() || isTyping.value) return
  
  const userMessage = currentPrompt.value.trim()
  
  // Add user message
  messages.value.push({
    type: 'user',
    text: userMessage,
    timestamp: new Date().toLocaleTimeString()
  })
  
  // Clear input
  currentPrompt.value = ''
  
  // Scroll to bottom
  await scrollToBottom()
  
  // Show typing indicator
  isTyping.value = true
  
  try {
    // Call your API here
    // const response = await yourApiCall(userMessage)
    const aiResponse = await sendToAI(userMessage)
    
    // Add AI response
    messages.value.push({
      type: 'bot',
      text: aiResponse,
      timestamp: new Date().toLocaleTimeString()
    })
  } catch (error) {
    // Handle error
    messages.value.push({
      type: 'bot',
      text: 'Sorry, I encountered an error. Please try again.',
      timestamp: new Date().toLocaleTimeString(),
      error: true
    })
  } finally {
    isTyping.value = false
    await scrollToBottom()
  }
}

// Add new line with Shift+Enter
const addNewLine = () => {
  currentPrompt.value += '\n'
}

// Format message (for code blocks, links, etc.)
const formatMessage = (text) => {
  // Convert URLs to links
  let formattedText = text.replace(/(https?:\/\/[^\s]+)/g, '<a href="$1" target="_blank">$1</a>')
  
  // Convert code blocks (simple version)
  formattedText = formattedText.replace(/`([^`]+)`/g, '<code>$1</code>')
  
  // Convert line breaks
  formattedText = formattedText.replace(/\n/g, '<br>')
  
  return formattedText
}

// Scroll to bottom of messages
const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// Start new chat
const startNewChat = () => {
  messages.value = []
  currentPrompt.value = ''
}

// Set suggested prompt
const setPrompt = (prompt) => {
  currentPrompt.value = prompt
  sendMessage()
}

// Auto-scroll on mount
onMounted(() => {
  scrollToBottom()
})
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: var(--surface-ground);
}

/* Header Styles */
.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: var(--surface-card);
  border-bottom: 1px solid var(--surface-border);
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.header-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-content h1 {
  margin: 0;
  font-size: 1.5rem;
  color: var(--text-color);
}

/* Messages Area */
.messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* Welcome Screen */
.welcome-screen {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: var(--text-color-secondary);
}

.welcome-screen h2 {
  margin: 1rem 0 0.5rem;
  font-size: 2rem;
}

.welcome-screen p {
  font-size: 1.1rem;
  margin-bottom: 2rem;
}

.suggestions {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 2rem;
}

.suggestion-chip {
  cursor: pointer;
  transition: transform 0.2s;
}

.suggestion-chip:hover {
  transform: translateY(-2px);
  background-color: var(--primary-color);
  color: white;
}

/* Message Styles */
.message {
  display: flex;
  gap: 1rem;
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message.user {
  justify-content: flex-end;
}

.message.user .message-content {
  background-color: var(--primary-color);
  color: white;
  border-radius: 1rem 1rem 0.25rem 1rem;
}

.message.bot .message-content {
  background-color: var(--surface-card);
  border: 1px solid var(--surface-border);
  border-radius: 1rem 1rem 1rem 0.25rem;
}

.message-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--surface-border);
}

.message.user .message-avatar {
  background-color: var(--primary-color);
  color: white;
  order: 2;
}

.message-content {
  max-width: 70%;
  padding: 0.75rem 1rem;
}

.message-text {
  line-height: 1.5;
  word-wrap: break-word;
}

.message-text code {
  background-color: rgba(0,0,0,0.1);
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.9em;
}

.message-time {
  font-size: 0.7rem;
  margin-top: 0.25rem;
  opacity: 0.7;
}

/* Typing Indicator */
.typing-indicator {
  display: flex;
  gap: 0.25rem;
  padding: 0.5rem 0;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: var(--text-color-secondary);
  animation: typing 1.4s infinite;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.4;
  }
  30% {
    transform: translateY(-10px);
    opacity: 1;
  }
}

/* Input Area */
.input-area {
  padding: 1rem 2rem;
  background-color: var(--surface-card);
  border-top: 1px solid var(--surface-border);
}

.input-container {
  display: flex;
  gap: 0.5rem;
  align-items: flex-end;
}

.chat-input {
  flex: 1;
  border-radius: 1.5rem;
  padding: 0.75rem 1rem;
  font-size: 1rem;
  resize: none;
}

.send-button {
  border-radius: 50%;
  width: 42px;
  height: 42px;
  padding: 0;
}

.input-footer {
  margin-top: 0.5rem;
  text-align: center;
  color: var(--text-color-secondary);
  font-size: 0.75rem;
}

/* Scrollbar Styling */
.messages-area::-webkit-scrollbar {
  width: 8px;
}

.messages-area::-webkit-scrollbar-track {
  background: var(--surface-border);
  border-radius: 4px;
}

.messages-area::-webkit-scrollbar-thumb {
  background: var(--primary-color);
  border-radius: 4px;
}

.messages-area::-webkit-scrollbar-thumb:hover {
  background: var(--primary-color-dark);
}

/* Responsive Design */
@media (max-width: 768px) {
  .chat-header {
    padding: 1rem;
  }
  
  .messages-area {
    padding: 1rem;
  }
  
  .message-content {
    max-width: 85%;
  }
  
  .input-area {
    padding: 1rem;
  }
  
  .suggestions {
    gap: 0.5rem;
  }
}
</style>