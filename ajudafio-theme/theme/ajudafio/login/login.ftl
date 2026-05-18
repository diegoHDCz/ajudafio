<#import "template.ftl" as layout>
<@layout.registrationLayout displayInfo=false displayWide=false; section>
    <#if section = "form">
        <div class="af-logo-wrap">
            <img src="${url.resourcesPath}/img/logo.png" alt="Ajuda Fio" class="af-form-logo" />
        </div>
        <h1 class="af-form-title">Bem-vindo de volta</h1>
        <p class="af-form-sub">Entre na sua conta para continuar</p>

        <#if realm.password>
            <form id="kc-form-login" action="${url.loginAction}" method="post">
                <div class="af-field">
                    <label for="username">
                        <#if !realm.loginWithEmailAllowed>Usuário<#elseif !realm.registrationEmailAsUsername>Usuário ou e-mail<#else>E-mail</#if>
                    </label>
                    <input tabindex="1" id="username" name="username"
                        value="${(login.username!'')}" type="text"
                        autofocus autocomplete="off"
                        placeholder="<#if !realm.loginWithEmailAllowed>seu usuário<#elseif !realm.registrationEmailAsUsername>usuário ou e-mail<#else>seu@email.com</#if>" />
                </div>

                <div class="af-field">
                    <label for="password">Senha</label>
                    <div class="af-input-wrap">
                        <input tabindex="2" id="password" name="password"
                            type="password" autocomplete="off" placeholder="••••••••" />
                    </div>
                </div>

                <div class="af-form-row">
                    <#if realm.rememberMe && !usernameEditDisabled??>
                        <label class="af-checkbox">
                            <input tabindex="3" id="rememberMe" name="rememberMe" type="checkbox"
                                <#if login.rememberMe??>checked</#if>>
                            Lembrar de mim
                        </label>
                    <#else>
                        <span></span>
                    </#if>
                    <#if realm.resetPasswordAllowed>
                        <a tabindex="5" href="${url.loginResetCredentialsUrl}" class="af-forgot">Esqueci minha senha</a>
                    </#if>
                </div>

                <input type="hidden" id="id-hidden-input" name="credentialId"
                    <#if auth.selectedCredential?has_content>value="${auth.selectedCredential}"</#if>/>

                <button tabindex="4" type="submit" class="af-btn" id="kc-login">
                    Entrar
                </button>
            </form>
        </#if>

        <#if realm.password && realm.registrationAllowed && !registrationDisabled??>
            <div class="af-register-link">
                <span>Não tem uma conta?</span>
                <a tabindex="6" href="${url.registrationUrl}">Criar conta</a>
            </div>
        </#if>
    </#if>
</@layout.registrationLayout>
