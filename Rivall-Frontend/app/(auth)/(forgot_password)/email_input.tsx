import { Box } from '@/components/ui/box';
import { VStack } from '@/components/ui/vstack';
import { Button, ButtonText, ButtonSpinner } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import { AlertCircleIcon } from '@/components/ui/icon';
import { Input, InputField } from "@/components/ui/input"
import {
  FormControl,
  FormControlError,
  FormControlErrorText,
  FormControlErrorIcon,
  FormControlLabel,
  FormControlLabelText,
} from "@/components/ui/form-control"
import { useState, useEffect } from "react"
import { sendCodeToEmail } from '@/api/auth';

import {
  validateEmail,
  validEmail,
} from '@/common/auth_helper_functions';

interface EmailInputProps {
  setStep: (step: string) => void;
  email: string;
  setEmail: (email: string) => void;
  codeSentState: string;
  setCodeSentState: (state: string) => void;
}

export function EmailInput({ setStep, email, setEmail, codeSentState, setCodeSentState }: EmailInputProps) {

  const [isInvalidEmail, setIsInvalidEmail] = useState(false)
  const [accountDoesntExist, setAccountDoesntExist] = useState(false)

  const handleSubmit = async () => {
    setCodeSentState('pending')

    // Validate Email
    if (!validateEmail(email)) {
      console.debug('Invalid Email')
      setIsInvalidEmail(true)
      return
    }

    // Send Password Recovery Code to Email
    const [data, success] = await sendCodeToEmail(email)
    if (success) {
      console.debug('Email Sent')
      setCodeSentState('sent')
      setStep('code')
      return
    }

    if (data === 'User not found') {
      console.debug('User not found')
      setAccountDoesntExist(true)
    }

    console.debug('Request Failed')
    setCodeSentState('failed')
  }

  useEffect(() => {
    setIsInvalidEmail(!validEmail(email))
    setAccountDoesntExist(false)
  }, [email])

  return (
    <VStack className="gap-4 mt-10">
      <Text className="text-typography-800 text-xl text-pretty text-left mb-3">Send us your email and we will send you a code to regain access to your account.</Text>
      <FormControl
        isInvalid={isInvalidEmail || accountDoesntExist}
        size="md"
        isDisabled={false}
        isReadOnly={false}
        isRequired={true}
        className="gap-4"
      >
        <Box>
          <FormControlLabel>
            <FormControlLabelText>Email</FormControlLabelText>
          </FormControlLabel>
          <Input className="my-1" size="md">
            <InputField
              type="text"
              placeholder="email"
              value={email}
              onChangeText={(text) => setEmail(text.toLowerCase())}
            />
          </Input>
          <FormControlError>
            <FormControlErrorIcon as={AlertCircleIcon} />
            <FormControlErrorText>
              {isInvalidEmail ? 'Invalid Email. ' : ''}
              {accountDoesntExist ? 'Account does not exist.' : ''}
            </FormControlErrorText>
          </FormControlError>
        </Box>
      </FormControl>
      <Button
        className="shadow-md shadow-black"
        size="lg"
        onPress={handleSubmit}
        disabled={codeSentState === 'pending'}
      >
        {codeSentState === 'pending' ? <ButtonSpinner/> : <ButtonText className="text-typography-0 text-lg">Send Code</ButtonText>}
      </Button>
    </VStack>
  )
}
