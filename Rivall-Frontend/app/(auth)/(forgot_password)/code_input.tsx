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

import { ResetPasswordState } from './password_reset'

export function VerificationCodeInput({ state }: { state: ResetPasswordState }) {

  const [isInvalidCode, setIsInvalidCode] = useState(false)

  const handleSubmit = async () => {
    state.setValidateCodeState('pending')

    // Validate Login logic
    if (!validateCode(state.code)) {
      console.debug('Invalid Email')
      return
    }

    // Send Login Request
    // const [data, success] = await login(email, password)
    // if (success) {
    //   console.debug('Login Successful')
    //   setUserData(data)
    //   router.replace('/entry')
    //   setLoginLoading(false)
    //   return
    // }

    state.setValidateCodeState('success')
    state.setStep('reset')
  }

  useEffect(() => {
    setIsInvalidCode(!validCode(state.code))
  }, [state.code])

  return (
    <VStack className="gap-4 mt-10">
      <Text className="text-typography-800 text-xl text-pretty text-left mb-3">A verification code has been sent to ‘{resetCtx.email}’.  Submit this code to reset your account password.</Text>
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
              placeholder="code"
              value={state.code}
              onChangeText={(text) => state.setCode(text.toLowerCase())}
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
        onPress={handleSubmit}
        disabled={state.validate_code === 'pending' || state.sent_email === 'pending'}
      >
        {state.validate_code === 'pending' ? <ButtonSpinner/> : <ButtonText className="text-typography-0 text-lg">Submit</ButtonText>}
      </Button>
      <Button
        className="shadow-md shadow-black bg-secondary-500"
        size="lg"
        variant="outline"
        onPress={handleSubmit}
        disabled={state.validate_code === 'pending' || state.sent_email === 'pending'}
      >
        {state.sent_email === 'pending' ? <ButtonSpinner/> : <ButtonText className="text-typography-800 text-lg">Resend Code</ButtonText>}
      </Button>
    </VStack>
  )
}
