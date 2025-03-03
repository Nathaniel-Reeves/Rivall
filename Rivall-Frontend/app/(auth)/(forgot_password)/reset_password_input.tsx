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

import { ResetPasswordState } from './password_reset'

export function ResetPasswordInput({ state }: { state: ResetPasswordState }) {

  const [passwordsMatch, setPasswordsMatch] = useState(true)
  const [isInvalidPassword, setIsInvalidPassword] = useState(false)

  const handleSubmit = async () => {
    state.setResetPasswordState('pending')

    // Validate Login logic
    if (!validatePassword(state.password)) {
      setIsInvalidPassword(true)
      return
    }

    if (!validatePassword(state.confirm_password)) {
      setIsInvalidPassword(true)
      return
    }

    if (!matchPassword(state.password, state.confirm_password)) {
      setPasswordsMatch(true)
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

    state.setResetPasswordState('success')
    // state.setStep('email')
  }

  const [showPassword, setShowPassword] = useState(false)
  const handleState = () => {
    setShowPassword((showState) => {
      return !showState
    })
  }

  useEffect(() => {
    setPasswordsMatch(matchPassword(state.password, state.confirm_password))
    setIsInvalidPassword(!validPassword(state.password) || !validPassword(state.confirm_password))
  }, [state.password, state.confirm_password])

  return (
    <VStack className="gap-4 mt-10">
      <Text className="text-typography-800 text-xl text-pretty text-left mb-3">Create a new password and regain access to your account.</Text>
      <FormControl
        isInvalid={!passwordsMatch || isInvalidPassword || state.reset_password == 'error'}
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
              value={state.password}
              onChangeText={(text) => state.setPassword(text)}
            />
            <InputSlot className="pr-3" onPress={handleState}>
              <InputIcon as={showPassword ? EyeIcon : EyeOffIcon} />
            </InputSlot>
          </Input>
          <FormControlError>
            <FormControlErrorIcon as={AlertCircleIcon} />
            <FormControlErrorText>
              Invalid Password.
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
              value={state.confirm_password}
              onChangeText={(text) => state.setConfirmPassword(text)}
            />
            <InputSlot className="pr-3" onPress={handleState}>
              <InputIcon as={showPassword ? EyeIcon : EyeOffIcon} />
            </InputSlot>
          </Input>
          <FormControlError>
            <FormControlErrorIcon as={AlertCircleIcon} />
            <FormControlErrorText>
              Invalid Password.
            </FormControlErrorText>
          </FormControlError>
        </Box>
      </FormControl>
      <Button
        className="shadow-md shadow-black"
        size="lg"
        onPress={handleSubmit}
        disabled={state.reset_password == 'pending'}
      >
        {state.reset_password == 'pending' ? <ButtonSpinner/> : <ButtonText className="text-typography-0 text-lg">Reset Password</ButtonText>}
      </Button>
    </VStack>
  )
}
