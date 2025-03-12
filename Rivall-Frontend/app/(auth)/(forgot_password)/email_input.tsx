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
import { useRouter } from 'expo-router';
import { sendCodeToEmail } from '@/api/auth';

import {
  validateEmail,
  validEmail,
} from '@/common/auth_helper_functions';

interface EmailInputProps {
  email: string;
  setEmail: (email: string) => void;
}

export function EmailInput({ email, setEmail }: EmailInputProps) {

  const [isInvalidEmail, setIsInvalidEmail] = useState(false)
  const [loading, setLoading] = useState(false)

  const router = useRouter()

  const handleSubmit = async () => {
    setLoading(true)

    // Validate Login logic
    if (!validateEmail(email)) {
      console.debug('Invalid Email')
      setIsInvalidEmail(true)
      return
    }

    // Send Password Recovery Code to Email
    const [data, success] = await sendCodeToEmail(email)
    if (success) {
      console.debug('Login Successful')
      setLoading(false)
      router.replace('/entry')
      return
    }

    console.debug('Request Failed')
    console.debug(data)
    console.debug(success)
    setLoading(false)
  }

  useEffect(() => {
    setIsInvalidEmail(!validEmail(email))
  }, [email])

  return (
    <VStack className="gap-4 mt-10">
      <Text className="text-typography-800 text-xl text-pretty text-left mb-3">Send us your email and we will send you a code to regain access to your account.</Text>
      <FormControl
        isInvalid={isInvalidEmail}
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
              Invalid Email.
            </FormControlErrorText>
          </FormControlError>
        </Box>
      </FormControl>
      <Button
        className="shadow-md shadow-black"
        size="lg"
        onPress={handleSubmit}
        disabled={loading}
      >
        {loading ? <ButtonSpinner/> : <ButtonText className="text-typography-0 text-lg">Send Code</ButtonText>}
      </Button>
    </VStack>
  )
}
