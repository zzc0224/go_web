import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import Content from '../views/Content.vue'
import Publish from '../views/Publish.vue'
import Login from '../views/Login.vue'
import SignUp from '../views/SignUp.vue'
import User from '../views/User.vue'
import OtherUser from "../views/OtherUser";
import Search from "../views/Search";
import Community from "@/views/Community";
const originalPush = VueRouter.prototype.push;
VueRouter.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err);
}
Vue.use(VueRouter)

  const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/post/:id',
    name: 'Content',
    component: Content
  },
  {
    path: '/publish',
    name: 'Publish',
    component: Publish,
    meta: { requireAuth: true }
  },
  {
    path: '/login',
    name:"Login",
    component: Login
  },
  {
    path: '/signup',
    name:"SignUp",
    component: SignUp
  },
  {
    path:'/user',
    name:"User",
    component: User
  },
  {
    path:'/otherUser/:userId',
    name:"OtherUser",
    component: OtherUser
  },
  {
    path: '/search/:keywords',
    name:"Search",
    component: Search
  },
    {
      path:`/community`,
      name:"Community",
      component: Community
    }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
