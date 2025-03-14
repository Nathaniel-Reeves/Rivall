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
import {
  validateCode,
  validCode,
} from '@/common/auth_helper_functions';
import { sendCodeToEmail, validateAccountRecoveryCode } from '@/api/auth';
import { useUserStore } from '@/global-store/user_store';

interface VerificationCodeInputProps {
  setStep: (step: string) => void;
  code: string;
  setCode: (code: string) => void;
  codeValidatedState: string;
  setCodeValidatedState: (state: string) => void;
  email: string;
  codeSentState: string;
  setCodeSentState: (state: string) => void;
}

export function VerificationCodeInput({ setStep, code, setCode, codeValidatedState, setCodeValidatedState, email, codeSentState, setCodeSentState }: VerificationCodeInputProps) {

  const [isInvalidCode, setIsInvalidCode] = useState(false)
  const [codeResent, setCodeResent] = useState(false)
  const setStoreData = useUserStore((state: any) => state.setStoreData)

  const handleRequestNewCode = async () => {
    setCodeSentState('pending')

    // Send Password Recovery Code to Email
    const [data, success] = await sendCodeToEmail(email)
    if (success) {
      console.debug('Email Sent')
      setCodeSentState('sent')
      setCodeResent(true)
      return
    }

    console.debug('Request Failed')
    console.debug(data)
    console.debug(success)
  }

  const handleSubmitCode = async () => {
    setCodeValidatedState('pending')
    setIsInvalidCode(false)

    // Validate Code
    if (!validateCode(code)) {
      console.debug('Invalid Code')
      setIsInvalidCode(true)
      setCodeValidatedState('none')
      return
    }

    // Test Code
    const [data, success] = await validateAccountRecoveryCode(email, code)
    if (success) {
      console.debug('Code Valid')
      setCodeValidatedState('validated')
      setStoreData(data)
      setStep('reset')
      return
    } else {
      if (data === 'Invalid code') {
        console.debug('Invalid Code')
        setIsInvalidCode(true)
      }
    }

    console.debug('Code Failed')
    setCodeValidatedState('failed')
  }

  useEffect(() => {
    setIsInvalidCode(!validCode(code))
  }, [])

  return (
    <VStack className="gap-4 mt-10">
      <Text className="text-typography-800 text-xl text-pretty text-left mb-3">A verification code has been sent to ‘{email}’.  Submit this code to reset your account password.</Text>
      {codeResent ? <Text className="text-typography-800 text-lg text-pretty text-left mb-3">Code resent to ‘{email}’.</Text> : null}
      <FormControl
        isInvalid={isInvalidCode}
        size="md"
        isDisabled={false}
        isReadOnly={false}
        isRequired={true}
        className="gap-4"
      >
        <Box>
          <FormControlLabel>
            <FormControlLabelText>Verification Code</FormControlLabelText>
          </FormControlLabel>
          <Input className="my-1" size="md">
            <InputField
              type="text"
              value={code}
              onChangeText={(text) => setCode(text.toUpperCase())}
              className="text-center tracking-[2em] text-xl"
            />
          </Input>
          <FormControlError>
            <FormControlErrorIcon as={AlertCircleIcon} />
            <FormControlErrorText>
              Invalid Code.
            </FormControlErrorText>
          </FormControlError>
        </Box>
      </FormControl>
      <Button
        className="shadow-md shadow-black"
        size="lg"
        onPress={handleSubmitCode}
        disabled={codeValidatedState === 'pending' || codeSentState === 'pending'}
      >
        {codeValidatedState === 'pending' || codeSentState === 'pending' ? <ButtonSpinner/> : <ButtonText className="text-typography-0 text-lg">Submit</ButtonText>}
      </Button>
      <Button
        className="shadow-md shadow-black bg-secondary-500"
        size="lg"
        variant="outline"
        onPress={handleRequestNewCode}
        disabled={codeValidatedState === 'pending' || codeSentState === 'pending'}
      >
        {codeValidatedState === 'pending' || codeSentState === 'pending' ? <ButtonSpinner/> : <ButtonText className="text-typography-800 text-lg">Resend Code</ButtonText>}
      </Button>
    </VStack>
  )
}
