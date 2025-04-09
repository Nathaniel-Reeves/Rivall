import { Box } from '@/components/ui/box';
import { VStack } from '@/components/ui/vstack';
import { Button, ButtonText, ButtonSpinner } from '@/components/ui/button';
import { Text } from '@/components/ui/text';
import { AlertCircleIcon } from '@/components/ui/icon';
import { Input, InputField, InputSlot, InputIcon } from "@/components/ui/input"
import {
  FormControl,
  FormControlError,
  FormControlErrorText,
  FormControlErrorIcon,
  FormControlLabel,
  FormControlLabelText,
} from "@/components/ui/form-control"
import { EyeIcon, EyeOffIcon } from "@/components/ui/icon"
import { useState, useEffect } from "react"
import {
  matchPassword,
  validatePassword,
  validPassword
} from '@/common/auth_helper_functions';
import { resetPassword } from '@/api/auth';
import { useUserStore } from '@/global-store/user_store';
import { useRouter } from 'expo-router';

interface ResetPasswordInputProps {
  setStep: (step: string) => void;
  password: string;
  setPassword: (password: string) => void;
  resetPasswordState: string;
  setResetPasswordState: (state: string) => void;
}

export function ResetPasswordInput({ password, setPassword, resetPasswordState, setResetPasswordState }: ResetPasswordInputProps) {

  const [confirm_password, setConfirmPassword] = useState('')
  const [passwordsMatch, setPasswordsMatch] = useState(true)
  const [isInvalidPassword, setIsInvalidPassword] = useState(false)
  
  const state = useUserStore((state: any) => state)
  const router = useRouter()

  const handleSubmit = async () => {
    setResetPasswordState('pending')

    // Validate Login logic
    if (!validatePassword(password)) {
      setIsInvalidPassword(true)
      return
    }

    if (!validatePassword(confirm_password)) {
      setIsInvalidPassword(true)
      return
    }

    if (!matchPassword(password, confirm_password)) {
      setPasswordsMatch(false)
      return
    }

    // Send Reset Password Request
    const res = await resetPassword(state, password)
    if (res === null) {
      console.debug('Password Reset Failed')
      setResetPasswordState('error')
      return
    }

    console.debug('Password Reset Successful')
    setResetPasswordState('success')
    router.replace('/entry/home')
  }

  const [showPassword, setShowPassword] = useState(false)
  const handleState = () => {
    setShowPassword((showState) => {
      return !showState
    })
  }

  useEffect(() => {
    setPasswordsMatch(matchPassword(password, confirm_password))
    setIsInvalidPassword(!validPassword(password) || !validPassword(confirm_password))
  }, [password, confirm_password])

  return (
    <VStack className="gap-4 mt-10">
      <Text className="text-typography-800 text-xl text-pretty text-left mb-3">Create a new password and regain access to your account.</Text>
      <FormControl
        isInvalid={!passwordsMatch || isInvalidPassword || resetPasswordState == 'error'}
        size="md"
        isDisabled={false}
        isReadOnly={false}
        isRequired={true}
        className="gap-4 mb-4"
      >
        <Box>
          <FormControlLabel>
            <FormControlLabelText>Password</FormControlLabelText>
          </FormControlLabel>
          <Input className="my-1" size="md">
            <InputField
              type={showPassword ? "text" : "password"}
              placeholder="password"
              value={password}
              onChangeText={(text) => setPassword(text)}
            />
            <InputSlot className="pr-3" onPress={handleState}>
              <InputIcon as={showPassword ? EyeIcon : EyeOffIcon} />
            </InputSlot>
          </Input>
          <FormControlError>
            <FormControlErrorIcon as={AlertCircleIcon} />
            <FormControlErrorText>
              {isInvalidPassword ? "Invalid Password.  " : ""}
              {passwordsMatch ? "" : "Passwords do not match."}
              {resetPasswordState == 'error' ? "An error occurred.  Please try again." : ""}
            </FormControlErrorText>
          </FormControlError>
        </Box>
        <Box>
          <FormControlLabel>
            <FormControlLabelText>Confirm Password</FormControlLabelText>
          </FormControlLabel>
          <Input className="my-1" size="md">
            <InputField
              type={showPassword ? "text" : "password"}
              placeholder="confirm password"
              value={confirm_password}
              onChangeText={(text) => setConfirmPassword(text)}
            />
            <InputSlot className="pr-3" onPress={handleState}>
              <InputIcon as={showPassword ? EyeIcon : EyeOffIcon} />
            </InputSlot>
          </Input>
          <FormControlError>
            <FormControlErrorIcon as={AlertCircleIcon} />
            <FormControlErrorText>
              {isInvalidPassword ? "Invalid Password.  " : ""}
              {passwordsMatch ? "" : "Passwords do not match."}
              {resetPasswordState == 'error' ? "An error occurred.  Please try again." : ""}
            </FormControlErrorText>
          </FormControlError>
        </Box>
      </FormControl>
      <Button
        className="shadow-md shadow-black"
        size="lg"
        onPress={handleSubmit}
        disabled={resetPasswordState == 'pending'}
      >
        {resetPasswordState == 'pending' ? <ButtonSpinner/> : <ButtonText className="text-typography-0 text-lg">Reset Password</ButtonText>}
      </Button>
    </VStack>
  )
}
