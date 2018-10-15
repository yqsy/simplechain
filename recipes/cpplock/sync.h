//
// Created by yq on 18-10-15.
//

#ifndef CPPLOCK_SYNC_H
#define CPPLOCK_SYNC_H

#include <condition_variable>
#include <thread>
#include <mutex>

#ifdef __clang__
// TL;DR Add GUARDED_BY(mutex) to member variables. The others are
// rarely necessary. Ex: int nFoo GUARDED_BY(cs_foo);
//
// See http://clang.llvm.org/docs/LanguageExtensions.html#threadsafety
// for documentation.  The clang compiler can do advanced static analysis
// of locking when given the -Wthread-safety option.
#define LOCKABLE __attribute__((lockable))
#define SCOPED_LOCKABLE __attribute__((scoped_lockable))
#define GUARDED_BY(x) __attribute__((guarded_by(x)))
#define GUARDED_VAR __attribute__((guarded_var))
#define PT_GUARDED_BY(x) __attribute__((pt_guarded_by(x)))
#define PT_GUARDED_VAR __attribute__((pt_guarded_var))
#define ACQUIRED_AFTER(...) __attribute__((acquired_after(__VA_ARGS__)))
#define ACQUIRED_BEFORE(...) __attribute__((acquired_before(__VA_ARGS__)))
#define EXCLUSIVE_LOCK_FUNCTION(...) __attribute__((exclusive_lock_function(__VA_ARGS__)))
#define SHARED_LOCK_FUNCTION(...) __attribute__((shared_lock_function(__VA_ARGS__)))
#define EXCLUSIVE_TRYLOCK_FUNCTION(...) __attribute__((exclusive_trylock_function(__VA_ARGS__)))
#define SHARED_TRYLOCK_FUNCTION(...) __attribute__((shared_trylock_function(__VA_ARGS__)))
#define UNLOCK_FUNCTION(...) __attribute__((unlock_function(__VA_ARGS__)))
#define LOCK_RETURNED(x) __attribute__((lock_returned(x)))
#define LOCKS_EXCLUDED(...) __attribute__((locks_excluded(__VA_ARGS__)))
#define EXCLUSIVE_LOCKS_REQUIRED(...) __attribute__((exclusive_locks_required(__VA_ARGS__)))
#define SHARED_LOCKS_REQUIRED(...) __attribute__((shared_locks_required(__VA_ARGS__)))
#define NO_THREAD_SAFETY_ANALYSIS __attribute__((no_thread_safety_analysis))
#define ASSERT_EXCLUSIVE_LOCK(...) __attribute((assert_exclusive_lock(__VA_ARGS__)))
#else
#define LOCKABLE
#define SCOPED_LOCKABLE
#define GUARDED_BY(x)
#define GUARDED_VAR
#define PT_GUARDED_BY(x)
#define PT_GUARDED_VAR
#define ACQUIRED_AFTER(...)
#define ACQUIRED_BEFORE(...)
#define EXCLUSIVE_LOCK_FUNCTION(...)
#define SHARED_LOCK_FUNCTION(...)
#define EXCLUSIVE_TRYLOCK_FUNCTION(...)
#define SHARED_TRYLOCK_FUNCTION(...)
#define UNLOCK_FUNCTION(...)
#define LOCK_RETURNED(x)
#define LOCKS_EXCLUDED(...)
#define EXCLUSIVE_LOCKS_REQUIRED(...)
#define SHARED_LOCKS_REQUIRED(...)
#define NO_THREAD_SAFETY_ANALYSIS
#define ASSERT_EXCLUSIVE_LOCK(...)
#endif


template<typename PARENT>
class LOCKABLE AnnotatedMixin : public PARENT {
public:
    void lock() EXCLUSIVE_LOCK_FUNCTION() {
        PARENT::lock();
    }

    void unlock() UNLOCK_FUNCTION() {
        PARENT::unlock();
    }

    bool try_lock()  EXCLUSIVE_TRYLOCK_FUNCTION(true) {
        return PARENT::try_lock();
    }
};

class CCriticalSection : public AnnotatedMixin<std::recursive_mutex> {
public:
    ~CCriticalSection() {
        // nothing
    }
};


class SCOPED_LOCKABLE CCriticalBlock {
private:
    std::unique_lock<CCriticalSection> lock;

    void Enter(const char *pszName, const char *pszFile, int nLine) {
        lock.lock();
    }

    bool TryEnter(const char *pszName, const char *pszFile, int nLine) {
        lock.try_lock();
        return lock.owns_lock();
    }

public:
    CCriticalBlock(CCriticalSection &mutexIn, const char *pszName, const char *pszFile, int nLine,
                   bool fTry = false)  EXCLUSIVE_LOCK_FUNCTION(mutexIn) : lock(mutexIn, std::defer_lock) {
        if (fTry) {
            TryEnter(pszName, pszFile, nLine);
        } else {
            Enter(pszName, pszFile, nLine);
        }
    }

    CCriticalBlock(CCriticalSection *pmutexIn, const char *pszName, const char *pszFile, int nLine,
                   bool fTry = false) EXCLUSIVE_LOCK_FUNCTION(pmutexIn) {
        if (!pmutexIn) {
            return;
        }

        lock = std::unique_lock<CCriticalSection>(*pmutexIn, std::defer_lock);
        if (fTry) {
            TryEnter(pszName, pszFile, nLine);
        } else {
            Enter(pszName, pszFile, nLine);
        }
    }

    ~CCriticalBlock() {

    }

    operator bool() {
        return lock.owns_lock();
    }
};

#define PASTE(x, y) x ## y
#define PASTE2(x, y) PASTE(x, y)

// 对象名加上__COUNTER__数字,每次加锁对象的名称都不同
// 1. cs 为外部传入的CCriticalSection对象
// 2. #cs 为外部传入的对象的名称,只做debug使用
// 3. __FILE__ 为文件名称,只做debug使用
// 4. __LINE__ 为行号,只做debug使用
#define LOCK(cs) CCriticalBlock PASTE2(criticalblock,__COUNTER__)(cs, #cs, __FILE__, __LINE__)

// 和上面的相同,只不过锁两个对象
#define LOCK2(cs1, cs2) CCriticalBlock criticalblock1(cs1, #cs1, __FILE__, __LINE__), criticalblock2(cs2, #cs2, __FILE__, __LINE__)


// 对象名为自己定义
// 1. cs 为外部传入的CCriticalSection对象, 只用作临时的尝试解锁
// 2. #cs 为外部传的对象的名称...
...
#define TRY_LOCK(cs,name) CCriticalBlock name(cs, #cs, __FILE__, __LINE__, true)

#endif //CPPLOCK_SYNC_H
