import { generateText } from "ai";
import { openai } from "@ai-sdk/openai";

export async function triageSupportTicket(ticketContent: string, user_email: string) {
  const result = await generateText({
    model: openai("gpt-4o-mini"),
    system: "You are an internal IT helpdesk categorization bot.",
    prompt: `Categorize the following support query for user ${user_email}: ${ticketContent}`,
  });

  return result.text;
}
