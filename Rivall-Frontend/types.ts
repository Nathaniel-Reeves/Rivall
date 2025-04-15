export type Message = {
  _id: string;
  user_id: string;
  timestamp: string;
  message_type: string;
  message_data: string;
  seen_by: string[];
};

export type User = {
  _id: string;
  first_name: string;
  last_name: string;
  email: string;
  avatar_image: string;
}

export type Contact = {
  _id: string;
  group_name: string;
  num_members: number;
  is_challenge: boolean;
  challenge_start_date: string;
  unread_count: number;
  last_message: Message;
  unread: boolean;
}

export type ChatMembers = {
  [key: string]: User;
}

export type GroupChallenge = {
  _id: string;
  group_name: string;
  num_members: number;
  chat_members: ChatMembers;
  start_date: string;
  last_message: Message;
  unread: boolean;
  unread_count: number;
}