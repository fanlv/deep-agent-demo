import { AgentEvent } from '../types';

export interface SSEClientOptions {
  url: string;
  onEvent: (event: AgentEvent) => void;
  onError?: (error: Error) => void;
  onComplete?: () => void;
}

export class SSEClient {
  private eventSource: EventSource | null = null;
  private abortController: AbortController | null = null;

  async connect(options: SSEClientOptions & { body?: string }): Promise<void> {
    const { url, onEvent, onError, onComplete, body } = options;

    this.abortController = new AbortController();

    try {
      console.log('[SSEClient] fetching:', url, 'method:', body ? 'POST' : 'GET');
      const response = await fetch(url, {
        method: body ? 'POST' : 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'text/event-stream',
        },
        body,
        signal: this.abortController.signal,
      });

      if (!response.ok) {
        const errBody = await response.text().catch(() => '');
        let errMsg = `HTTP ${response.status}`;
        try {
          const errJson = JSON.parse(errBody);
          if (errJson.error) errMsg = errJson.error;
        } catch {
          if (errBody) errMsg = errBody;
        }
        throw new Error(errMsg);
      }
      console.log('[SSEClient] response ok, reading stream...');

      const reader = response.body?.getReader();
      if (!reader) {
        throw new Error('No readable stream');
      }

      const decoder = new TextDecoder();
      let buffer = '';

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });
        const events = buffer.split('\n\n');
        buffer = events.pop() || '';

        for (const eventBlock of events) {
          if (!eventBlock.trim()) continue;
          
          const lines = eventBlock.split('\n');
          let data = '';
          
          for (const line of lines) {
            if (line.startsWith('data:')) {
              data = line.slice(5).trim();
            }
          }
          
          if (!data) continue;
          
          if (data === '[DONE]') {
            onComplete?.();
            return;
          }
          
          try {
            const event = JSON.parse(data) as AgentEvent;
            onEvent(event);
          } catch {
            console.warn('Failed to parse SSE data:', data);
          }
        }
      }

      onComplete?.();
    } catch (error) {
      if ((error as Error).name !== 'AbortError') {
        onError?.(error as Error);
      }
    }
  }

  disconnect(): void {
    if (this.abortController) {
      this.abortController.abort();
      this.abortController = null;
    }
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
  }
}
