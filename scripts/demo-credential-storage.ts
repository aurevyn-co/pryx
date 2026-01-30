#!/usr/bin/env bun
/**
 * Credential Storage Demonstration
 * 
 * This script demonstrates how Pryx stores credentials securely using:
 * - Argon2id key derivation
 * - AES-256-GCM encryption
 * - Atomic file writes
 * - Secure memory management
 * 
 * NOTE: This uses TEST data only - never commit real credentials!
 */

import { PasswordManager } from '../packages/vault/src/password-manager.js';
import { VaultStorage } from '../packages/vault/src/storage.js';
import { mkdtempSync, rmSync } from 'fs';
import { tmpdir } from 'os';
import { join } from 'path';

// Test data - NEVER use real credentials in tests!
const TEST_MASTER_PASSWORD = 'demo-master-password-do-not-use-in-production';
const TEST_SERVICE_NAME = 'z.ai';
const TEST_API_KEY_FORMAT = 'sk-demo-1234567890abcdef-TEST-ONLY';

async function demonstrateCredentialStorage() {
  console.log('🔐 Pryx Credential Storage Demonstration\n');
  
  // Create temporary directory for demo
  const tempDir = mkdtempSync(join(tmpdir(), 'pryx-demo-'));
  const vaultPath = join(tempDir, 'credentials.vault');
  
  try {
    // Step 1: Initialize Password Manager
    console.log('Step 1: Initializing Password Manager...');
    const passwordManager = new PasswordManager({
      autoLockMs: 5 * 60 * 1000, // 5 minutes
    });
    
    // Step 2: Unlock vault with master password
    console.log('Step 2: Unlocking vault with master password...');
    await passwordManager.unlock(TEST_MASTER_PASSWORD);
    console.log('   ✅ Vault unlocked successfully\n');
    
    // Step 3: Prepare credential data
    console.log('Step 3: Preparing credential data...');
    const credentialData = {
      service: TEST_SERVICE_NAME,
      apiKey: TEST_API_KEY_FORMAT, // In real usage, this would be the actual key
      metadata: {
        createdAt: new Date().toISOString(),
        provider: 'z.ai',
        type: 'api-key',
      },
    };
    
    const plaintext = Buffer.from(JSON.stringify(credentialData));
    console.log(`   📄 Data size: ${plaintext.length} bytes`);
    console.log(`   🔑 Service: ${credentialData.service}`);
    console.log(`   📝 Key format: ${TEST_API_KEY_FORMAT.substring(0, 10)}...${TEST_API_KEY_FORMAT.substring(TEST_API_KEY_FORMAT.length - 8)}\n`);
    
    // Step 4: Encrypt the credential
    console.log('Step 4: Encrypting credential with AES-256-GCM...');
    const encrypted = await passwordManager.encrypt(plaintext);
    console.log(`   🔒 Encrypted size: ${encrypted.ciphertext.length} bytes`);
    console.log(`   🧂 Salt: ${encrypted.salt.toString('hex').substring(0, 16)}...`);
    console.log(`   🔢 IV: ${encrypted.iv.toString('hex')}`);
    console.log(`   🏷️  Auth Tag: ${encrypted.tag.toString('hex').substring(0, 16)}...\n`);
    
    // Step 5: Store in vault file
    console.log('Step 5: Storing in encrypted vault file...');
    const storage = new VaultStorage();
    const vaultFile = await storage.load(vaultPath, TEST_MASTER_PASSWORD).catch(() => ({
      version: 1,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      metadata: {
        salt: encrypted.salt.toString('base64'),
        algorithm: 'argon2id+aes-256-gcm',
        iterations: 3,
        memoryCost: 65536,
      },
      entries: [],
    }));
    
    // Add entry to vault
    const entry = await storage.addEntry(
      vaultFile,
      {
        type: 'api-key',
        name: 'z.ai API Key',
        data: credentialData,
      },
      TEST_MASTER_PASSWORD
    );
    
    console.log(`   💾 Entry ID: ${entry.id}`);
    console.log(`   📁 Vault file: ${vaultPath}`);
    console.log(`   🔐 File permissions: 0o600 (user read/write only)\n`);
    
    // Step 6: Save vault atomically
    console.log('Step 6: Saving vault (atomic write with backup)...');
    await storage.save(vaultPath, vaultFile, TEST_MASTER_PASSWORD);
    console.log('   ✅ Vault saved successfully\n');
    
    // Step 7: Verify integrity
    console.log('Step 7: Verifying vault integrity...');
    const integrity = await storage.verifyIntegrity(vaultPath);
    console.log(`   ✅ Integrity check: ${integrity.valid ? 'PASSED' : 'FAILED'}`);
    console.log(`   📊 Entries: ${integrity.entryCount}`);
    console.log(`   🔍 Corrupted entries: ${integrity.corruptedEntries}\n`);
    
    // Step 8: Retrieve and decrypt
    console.log('Step 8: Retrieving and decrypting credential...');
    const retrievedEntry = await storage.getEntry(vaultFile, entry.id, TEST_MASTER_PASSWORD);
    const retrievedData = JSON.parse(Buffer.from(retrievedEntry.data as Buffer).toString());
    
    console.log(`   📋 Retrieved service: ${retrievedData.service}`);
    console.log(`   ✅ Data integrity verified: ${retrievedData.apiKey === credentialData.apiKey}\n`);
    
    // Step 9: Lock vault
    console.log('Step 9: Locking vault...');
    passwordManager.lock();
    console.log('   🔒 Vault locked. Key cleared from memory.\n');
    
    // Summary
    console.log('📊 Security Summary:');
    console.log('   • Master password: Never stored, only used for key derivation');
    console.log('   • Encryption: AES-256-GCM with unique IV per entry');
    console.log('   • Key derivation: Argon2id (64MB memory, 3 iterations)');
    console.log('   • File storage: Atomic writes with automatic backups');
    console.log('   • Memory: Secure clearing of sensitive data');
    console.log('   • Permissions: User read/write only (0o600)\n');
    
    console.log('✅ Demonstration complete!');
    
  } catch (error) {
    console.error('❌ Error:', error.message);
    process.exit(1);
  } finally {
    // Cleanup
    rmSync(tempDir, { recursive: true, force: true });
    console.log(`\n🧹 Cleaned up temporary files`);
  }
}

// Run demonstration
demonstrateCredentialStorage();
