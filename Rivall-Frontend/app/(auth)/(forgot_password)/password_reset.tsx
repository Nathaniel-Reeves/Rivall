import { Box } from '@/components/ui/box';
import { Card } from '@/components/ui/card';
import { Spinner } from "@/components/ui/spinner"
import { Button, ButtonIcon } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import { ChevronLeftIcon } from '@/components/ui/icon';
import { useState } from "react"
import { ScrollView } from 'react-native';
import { Link } from 'expo-router';
import EmailInput from './email_input';
import VerificationCodeInput from './code_input';
import ResetPasswordInput from './reset_password_input';

export default function ForgotPasswordScreens() {
  const [step, setStep] = useState('email')
  const [email, setEmail] = useState('')
  const [codeSentState, setCodeSentState] = useState('none')
  const [code, setCode] = useState('')
  const [codeValidatedState, setCodeValidatedState] = useState('none')
  const [resetPasswordState, setResetPasswordState] = useState('none')
  const [password, setPassword] = useState('')

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
      <Card className="align-top m-4 my-16 p-10 shadow-md shadow-black flex-col justify-between">
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
        {step === 'email' ? <EmailInput setStep={setStep} email={email} setEmail={setEmail} codeSentState={codeSentState} setCodeSentState={setCodeSentState}/> : null}
        {step === 'code' ? <VerificationCodeInput setStep={setStep} code={code} setCode={setCode} codeValidatedState={codeValidatedState} setCodeValidatedState={setCodeValidatedState} codeSentState={codeSentState} setCodeSentState={setCodeSentState} email={email}/> : null}
        {step === 'reset' ? <ResetPasswordInput setStep={setStep} password={password} setPassword={setPassword} resetPasswordState={resetPasswordState} setResetPasswordState={setResetPasswordState}/> : null}
      </Card>
    </ScrollView>
  );
}
