import { ref } from "vue";
import { CURRENT_USER_MOCK, type UserProfileMock } from "@/mocks/user.mock";

const isAuthenticated = ref(true);
const currentUser = ref<UserProfileMock>({ ...CURRENT_USER_MOCK });

export function useAuth() {
  function signIn(userData?: Partial<UserProfileMock>) {
    isAuthenticated.value = true;
    if (userData?.name) currentUser.value.name = userData.name;
    if (userData?.username) currentUser.value.username = userData.username;
  }

  function signOut() {
    isAuthenticated.value = false;
  }

  return {
    isAuthenticated,
    currentUser,
    signIn,
    signOut,
  };
}
