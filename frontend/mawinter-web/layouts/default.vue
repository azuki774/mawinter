<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div class="min-h-screen bg-slate-50">
    <header class="bg-white border-b border-slate-200 sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-14">
          <NuxtLink to="/" class="text-xl font-bold text-slate-800 hover:text-blue-600 transition-colors">
            mawinter
          </NuxtLink>
          <nav class="flex items-center gap-1">
            <NuxtLink
              to="/"
              class="px-3 py-2 rounded-md text-sm font-medium text-slate-600 hover:text-blue-600 hover:bg-blue-50 transition-colors"
              exact-active-class="!text-blue-600 !bg-blue-50"
            >
              トップ
            </NuxtLink>
            <NuxtLink
              to="/summary"
              class="px-3 py-2 rounded-md text-sm font-medium text-slate-600 hover:text-blue-600 hover:bg-blue-50 transition-colors"
              active-class="!text-blue-600 !bg-blue-50"
            >
              サマリー
            </NuxtLink>
            <NuxtLink
              to="/graph"
              class="px-3 py-2 rounded-md text-sm font-medium text-slate-600 hover:text-blue-600 hover:bg-blue-50 transition-colors"
              active-class="!text-blue-600 !bg-blue-50"
            >
              グラフ
            </NuxtLink>
          </nav>
        </div>
      </div>
    </header>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <slot />
    </main>
  </div>
</template>
