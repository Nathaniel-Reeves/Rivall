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

  const handleSubmit = async () => {
    setLoading(true)

    // Validate Login logic
    if (!validateEmail(email)) {
      console.debug('Invalid Email')
      setIsInvalidEmail(true)
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

    setTimeout(() => {
      setLoading(false)
    }, 2000)
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
