import { Box } from '@/components/ui/box';
import { Card } from '@/components/ui/card';
import { Spinner } from "@/components/ui/spinner"
import { Button, ButtonIcon } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import { ChevronLeftIcon } from '@/components/ui/icon';
import { useState } from "react"
import { ScrollView } from 'react-native';
import { Link } from 'expo-router';
import { EmailInput } from './email_input';
import { VerificationCodeInput } from './code_input';
import { ResetPasswordInput } from './reset_password_input';

export default function ForgotPasswordScreens() {
  const [step, setStep] = useState('email')
  const [email, setEmail] = useState('')
  const [codeSent, setCodeSent] = useState('none')
  const [code, setCode] = useState('')
  const [codeValidated, setCodeValidated] = useState('none')
  const [resetPasswordState, setResetPasswordState] = useState('none')

  return (
    <ScrollView
      style={{
        position: 'absolute',
        left: 0,
        right: 0,
        top: 0,
        height: '100%',
      }}
    >
      <Card className="align-top m-4 my-10 p-10 shadow-md shadow-black flex-col justify-between">
        <Box>
          <Link href="/welcome" asChild replace>
            <Button
              className="mb-6 w-20"
              variant="outline"
            >
              <ButtonIcon as={ChevronLeftIcon}/>
            </Button>
          </Link>
          <Text className="text-typography-800 text-4xl font-normal text-pretty text-left mb-4">
            {step === 'email' ? 'Forgot Password?' : ''}
            {step === 'code'  ? 'Check Your Inbox!' : ''}
            {step === 'reset' ? 'Welcome Back!' : ''}
          </Text>
          <Text className="text-typography-800 text-2xl text-pretty text-left mb-10">
            {step === 'email' ? 'Lets help you reset it!' : ''}
            {step === 'code'  ? 'Reply with your code!' : ''}
            {step === 'reset' ? 'Lets reset your password!' : ''}
          </Text>
        </Box>
        {step === 'email' ? <EmailInput email={email} setEmail={setEmail}/> : null}
        {/* {step === 'code' ? <VerificationCodeInput code={code} setEmail={setCode}/> : null}
        {step === 'reset' ? <ResetPasswordInput password={password} setEmail={setPassword}/> : null} */}
      </Card>
    </ScrollView>
  );
}
